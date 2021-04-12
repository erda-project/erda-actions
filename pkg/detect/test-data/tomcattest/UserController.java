package com.jjy.controller;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
@RequestMapping(value = "/user")
public class UserController{

  
	@RequestMapping("/list")
	@ResponseBody
	public String getApplist() {
		
		return "123gffdhgfdg";
	}
}
