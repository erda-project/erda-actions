package io.terminus.pampas.showcase.user.service.impl;

import io.terminus.boot.rpc.common.annotation.RpcProvider;
import io.terminus.pampas.showcase.user.model.User;
import io.terminus.pampas.showcase.user.service.UserService;
import lombok.val;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;


/**
 * Created by gengrong on 2017/4/18.
 */
@Service
@RpcProvider(version = "0.0.1")
public class UserServiceImpl implements UserService {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @Override
    public User login(String userName, String password) {
        if (StringUtils.isEmpty(userName) || StringUtils.isEmpty(password)) {
            return null;
        }
        val sql = "select * from pp_users where user_name = ?";
        User loginUser = jdbcTemplate.queryForObject(sql, new Object[]{userName}, User.class);
        if (loginUser != null && password.equals(loginUser.getPassword())) {
            return loginUser;
        }
        return null;
    }

    @Override
    public User getById(Long id) {
        val sql = "select * from pp_users where id = ?";
        User user = jdbcTemplate.queryForObject(sql, new Object[]{id}, new BeanPropertyRowMapper<User>(User.class));
        return user;
    }

    @Override
    public String healthCheck() {
        return "OK";
    }
}
