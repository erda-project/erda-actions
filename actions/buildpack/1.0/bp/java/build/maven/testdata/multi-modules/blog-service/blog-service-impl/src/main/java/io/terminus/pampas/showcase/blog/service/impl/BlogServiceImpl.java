package io.terminus.pampas.showcase.blog.service.impl;

import io.terminus.boot.rpc.common.annotation.RpcConsumer;
import io.terminus.boot.rpc.common.annotation.RpcProvider;
import io.terminus.pampas.showcase.blog.model.Blog;
import io.terminus.pampas.showcase.blog.service.BlogService;
import io.terminus.pampas.showcase.user.model.User;
import io.terminus.pampas.showcase.user.service.UserService;
import lombok.val;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.BeanPropertyRowMapper;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * Created by gengrong on 2017/4/18.
 */
@Service
@RpcProvider(version = "0.0.1")
public class BlogServiceImpl implements BlogService {

    @Autowired
    private JdbcTemplate jdbcTemplate;

    @RpcConsumer(version = "0.0.1", check = "false")
    private UserService userService;

    @Override
    public void createBlog(String title, String content) {
        val sql = "INSERT INTO pp_blogs (title, content) VALUES (?, ?)";
        jdbcTemplate.update(sql, new Object[]{title, content}, new Object[]{String.class, String.class});
    }

    @Override
    public Iterable<Blog> listAll() {
        val sql = "SELECT * FROM pp_blogs";
        List<Blog> blogs = jdbcTemplate.query(sql, new Object[]{}, new BeanPropertyRowMapper<Blog>(Blog.class));

        for (Blog blog : blogs) {
            User creator = userService.getById(blog.getCreator());
            if (creator != null) {
                blog.setCreatorName(creator.getNickName());
            }
        }
        return blogs;
    }

    @Override
    public Blog get(Long id) {
        val sql = "SELECT * FROM pp_blogs WHERE id = ?";
        Blog blog = jdbcTemplate.queryForObject(sql, new Object[]{id}, new BeanPropertyRowMapper<Blog>(Blog.class));
        User creator = userService.getById(blog.getCreator());
        if (creator != null) {
            blog.setCreatorName(creator.getNickName());
        }
        return blog;
    }
}
