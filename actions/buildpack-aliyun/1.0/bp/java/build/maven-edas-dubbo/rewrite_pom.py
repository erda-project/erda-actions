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
        d_groupId.text = "com.alibaba.edas"
        d_artifactId = ET.SubElement(dependency, "artifactId")
        d_artifactId.text = "edas-dubbo-extension"
        d_version = ET.SubElement(dependency, "version")
        d_version.text = "1.0.8"

        for c in root.iter(pre + "dependencies"):
            c.append(dependency)

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
        print("Rewrite %s successful for edas-dubbo." % args.pom_path)
    except:
        print("rewrite meet error")
        sys.exit(1)
