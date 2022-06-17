import { spawn } from 'child_process';
import color from 'ansi-colors';
import { writeFile } from 'fs/promises';


const { env } = process;
const NODE_VERSION = env.ACTION_NODE_VERSION || 14;
const PRESERVE_TIME = env.ACTION_PRESERVE_TIME;
const nodeVerMap = {
  12: 'v12.22.5',
  14: 'v14.19.0',
}
const faqUrl = 'https://docs.erda.cloud/latest/manual/faq/faq.html';

const logPrefix = '[js pack] '
const logInfo = (...msg) => console.log(color.greenBright(logPrefix + msg.join('')));
// const logSuccess = (...msg) => console.log(color.greenBright('✅', logPrefix));
const logWarn = (...msg) => console.log(color.yellowBright(logPrefix + msg.join('')));
const logError = (...msg) => console.error(color.redBright(logPrefix + msg.join('')));

const getBlockTitle = (title = '', char = '=') => {
  const half = 100 > title.length ? Math.floor((100 - title.length) / 2) : 0;
  return char.repeat(half) + ' ' + title + ' ' + char.repeat(half)
}


async function sleep(seconds) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, seconds * 1000);
  })
}

const metadata = [];
function output(name, value) {
  metadata.push({ name, value });
}

async function writeMetaFile() {
  if (metadata.length) {
    const data = JSON.stringify({ metadata });
    try {
      await writeFile(env.METAFILE, data)
      logInfo('write output to metafile success');
      logInfo(data);
    } catch (error) {
      logError('write output to metafile failed', error);
    }
  }
}

async function execCMD(cmd) {
  return new Promise(function (resolve, reject) {
    try {
      const childProcess = spawn(cmd, [], { stdio: [process.stdin, process.stdout, 'pipe'], shell: true, env: process.env });
      childProcess.on('close', resolve);
      childProcess.on('error', reject);
      childProcess.on('exit', (code, signal) => {
        if (code) {
          reject(`Command '${cmd}' failed with code ${code}`)
        } else if (signal) {
          reject(`Command '${cmd}' process was killed with signal ${signal}`);
        } else {
          // logInfo(`Command \`${cmd}\` execute finished`);
        }
      });
      childProcess.stderr.on('data', (data) => {
        const str = data.toString();
        process.stderr.write(data);

        const memoryOverflow = str.startsWith('npm ERR! errno 137');
        if (memoryOverflow) {
          logError(`Memory overflow, should set --max_old_space_size in build_cmd, view ${faqUrl} for detail`)
        }

        const happypackErr = str.startsWith('Happythread[babel:0] unable to send to worker');
        if (happypackErr) {
          logError(`Please set 'threads: 1' in happypack config, view ${faqUrl} for detail`)
        }

        // find npm error log file in output
        const npmLogFile = /\s+([\S]+-debug\.log)$/mg.exec(str);
        if (npmLogFile && npmLogFile[1]) {
          const fileName = npmLogFile[1];
          process.nextTick(() => { // prevent output out of order
            runCmd('cat ' + npmLogFile[1] + ' >&2', {
              before: () => logError(getBlockTitle(fileName + ' start')),
              after: () => logError(getBlockTitle(fileName + ' end'))
            });
          });
        }
      });

    } catch (error) {
      reject(error);
    }
  });
}

async function runCmd(cmdStr, hooks = {}) {
  if (!cmdStr) {
    return false;
  }

  if (hooks.before) {
    await hooks.before();
  }
  let success = false;
  await execCMD(cmdStr).catch(async e => {
    e && logError(e);

    if (PRESERVE_TIME) {
      logInfo('Job container preserve time: ', PRESERVE_TIME);
      logInfo(`You can use terminal to debug, view ${faqUrl} for detail`)
      await sleep(+PRESERVE_TIME)
    } else {
      logWarn('You can set \'preserve_time\' in action params to keep this job container running, and use terminal to debug');
    }
    success = false;
  }).finally(() => {
    if (hooks.after) hooks.after();
  });
  return success;
}

async function run() {
  logInfo(getBlockTitle('Build Env', '='));
  logInfo('Node Version: ' + nodeVerMap[NODE_VERSION]);
  // 流水线会把 WORKDIR 目录内的东西共享给后面的 action
  logInfo('Working directory: ', env.WORKDIR);
  logInfo('NAMESPACE: ', env.DICE_NAMESPACE);
  logInfo('PIPELINE_LIMITED_CPU: ', env.PIPELINE_LIMITED_CPU);
  logInfo('PIPELINE_LIMITED_MEM: ', env.PIPELINE_LIMITED_MEM);
  logInfo(getBlockTitle('Build Env', '='));
  if (!env.ACTION_WORKDIR) {
    logWarn('work_dir not set, generally you can set to "${{ dirs.git-checkout }}"');
    process.exit(1);
  }
  // 把代码复制到 WORKDIR 里，因为不知道编译输出的目录名是什么，没法只把编译完的内容复制过来
  await runCmd(`cp -r ${env.ACTION_WORKDIR}/* ${env.WORKDIR}`);
  logInfo('Copy git checkout files to working directory finished');

  const buildCmdPrefix = `. ~/.nvm/nvm.sh && nvm use ${NODE_VERSION} && `;
  const runResult = [];
  if (env.ACTION_BUILD_CMD && env.ACTION_BUILD_CMD.startsWith('["')) {
    const cmdList = JSON.parse(env.ACTION_BUILD_CMD);
    for (let i = 0; i < cmdList.length; i++) {
      const buildCmd = cmdList[i];
      logInfo(`Execute build_cmd part ${i + 1}: `, buildCmd)
      const runSuccess = await runCmd(buildCmdPrefix + buildCmd);
      runResult.push(runSuccess);
    }
  } else {
    const buildCmd = env.ACTION_BUILD_CMD || 'npm run build';
    logInfo('Execute build_cmd: ', buildCmd)
    const runSuccess = await runCmd(buildCmdPrefix + buildCmd);
    runResult.push(runSuccess);
  }
  await writeMetaFile();
  if (runResult.every(r => r !== false)) {
    logInfo('🎉 Build success')
  } else {
    setTimeout(() => {
      const fileMsg = `❌ Build failed, you can view ${faqUrl} to search for common failure causes`;
      // print to stdout and stderr both
      logWarn(fileMsg)
      logError(fileMsg)
      process.exit(1);
    }, 100)
  }
}

run();
