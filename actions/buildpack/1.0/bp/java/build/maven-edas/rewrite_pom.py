#!/usr/bin/python
# coding: utf-8

import xml.etree.cElementTree as ET
import re
import sys
import argparse

class ConfigXMLFile(object):

    def __init__(self, file):
        self.config = file  # 配置文件path
        self.tree = None

    def readXML(self, type):
        '''
        读取并解析xml文件
        return: ElementTree
        '''
        self.tree = ET.ElementTree()
        if type == "pom":
            XML_NS_NAME = ""
            XML_NS_VALUE = "http://maven.apache.org/POM/4.0.0"
            ET.register_namespace(XML_NS_NAME, XML_NS_VALUE)
        self.tree.parse(self.config)

    def writeXML(self, out_path):
        '''
        将xml文件写出
        out_path: 写出路径
        '''
        self.tree.write(out_path, encoding="utf-8", xml_declaration=True)

    def rewritePom(self, out_path):
        '''
        修改pom中的依赖包的version
        :param out_path: 修改后的配置文件路径
        :return:
        '''
        pre_sibling = None
        root = self.tree.getroot()  # 根node
        pre = (re.split('project', root.tag))[0]  # 获取pom元素tag的pre

        dependency = ET.Element("dependency")
        d_groupId = ET.SubElement(dependency, "groupId")
        d_groupId.text = "org.springframework.cloud"
        d_artifactId = ET.SubElement(dependency, "artifactId")
        d_artifactId.text = "spring-cloud-starter-pandora"
        d_version = ET.SubElement(dependency, "version")
        d_version.text = "1.2"

        for c in root.iter(pre + "dependencies"):
            c.append(dependency)

        remove_plugin = False
        for child in root.iter(pre + "plugins"):
            groupIdFound = False
            artifactIdFound = False
            build_del = None
            for p_child in child:
                for sub_child in p_child:
                    if sub_child.tag == (pre + "groupId") and sub_child.text == "org.springframework.boot":
                        groupIdFound = True
                    if sub_child.tag == (pre + "artifactId") and sub_child.text == "spring-boot-maven-plugin":
                        artifactIdFound = True

                    if groupIdFound and artifactIdFound and build_del == None:
                        build_del = p_child

            if build_del != None:
                child.remove(build_del)
                remove_plugin = True

        if remove_plugin:
            plugin = ET.Element("plugin")
            p_groupId = ET.SubElement(plugin, "groupId")
            p_groupId.text = "com.taobao.pandora"
            p_artifactId = ET.SubElement(plugin, "artifactId")
            p_artifactId.text = "pandora-boot-maven-plugin"
            p_version = ET.SubElement(plugin, "version")
            p_version.text = "2.1.7.8"
            p_executions = ET.SubElement(plugin, "executions")
            p_execution = ET.SubElement(p_executions, "execution")
            p_e_phase = ET.SubElement(p_execution, "phase")
            p_e_phase.text = "package"
            p_e_goals = ET.SubElement(p_execution, "goals")
            p_e_g_goal = ET.SubElement(p_e_goals, "goal")
            p_e_g_goal.text = "repackage"

            for p in root.iter(pre + "plugins"):
                p.append(plugin)

        self.writeXML(out_path)

if __name__ == "__main__":
    try:
        parser = argparse.ArgumentParser(description="rewrite pom.xml")
        parser.add_argument("pom_path")
        args = parser.parse_args()
        pom_xml = ConfigXMLFile(args.pom_path)
        pom_xml.readXML("pom")
        print("Read %s successful." % args.pom_path)
        pom_xml.rewritePom(args.pom_path)
        print("Rewrite %s successful for edas." % args.pom_path)
    except:
        print("rewrite meet error")
        sys.exit(1)
