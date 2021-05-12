package io.terminus.pampas.showcase.user.service;


import io.terminus.pampas.showcase.user.model.User;

public interface UserService {

    User login(String userName, String password);

    User getById(Long id);

    String healthCheck();
}
