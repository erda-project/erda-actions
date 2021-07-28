### 执行应用流水线

根据测试计划名称和 参数配置 来执行对应的测试用例，然后等待执行完成

__yaml

      - testplan-run:
          alias: testplan-run
          description: 根据自动化测试计划启动测试计划并等待完成
          version: "1.0"
          params:
            cms: autotest^scope-project-autotest-testcase^scopeid-2^344025938771587849
            test_plan: 10000093