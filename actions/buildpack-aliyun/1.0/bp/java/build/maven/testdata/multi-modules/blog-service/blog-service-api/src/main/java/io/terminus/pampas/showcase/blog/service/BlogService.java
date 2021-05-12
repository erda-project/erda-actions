package io.terminus.pampas.showcase.blog.service;


import io.terminus.pampas.showcase.blog.model.Blog;

public interface BlogService {

    void createBlog(String title, String content);

    Iterable<Blog> listAll();

    Blog get(Long id);
}
