package io.terminus.pampas.showcase.user.repository;


import io.terminus.pampas.showcase.user.model.User;

/**
 * Created by gengrong on 2017/4/18.
 */
public interface UserRepository {

    User findFirstByUserName(String userName);

}
