package io.terminus.pampas.showcase.restful;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.terminus.boot.rpc.common.annotation.RpcConsumer;
import io.terminus.pampas.showcase.user.model.User;
import io.terminus.pampas.showcase.user.service.UserService;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@Api
@RestController
@RequestMapping("/api/users")
public class Users {

    @RpcConsumer(version = "0.0.1", check = "false")
    private UserService userService;

    @ApiOperation("user login")
    @RequestMapping(value = "/login", method = RequestMethod.POST)
    User login(String userName, String password) {
        return userService.login(userName, password);
    }

}
