package io.terminus.pampas.showcase.restful;

import io.swagger.annotations.Api;
import io.swagger.annotations.ApiOperation;
import io.terminus.boot.rpc.common.annotation.RpcConsumer;
import io.terminus.pampas.showcase.blog.model.Blog;
import io.terminus.pampas.showcase.blog.service.BlogService;
import org.springframework.web.bind.annotation.*;

@Api
@RestController
@RequestMapping("/api/blogs")
public class Blogs {

    @RpcConsumer(version = "0.0.1", check = "false")
    private BlogService blogService;

    @ApiOperation("blog list")
    @RequestMapping(method = RequestMethod.GET)
    Iterable<Blog> list() {
        return blogService.listAll();
    }

    @ApiOperation("blog detail")
    @RequestMapping(value = "/{id}", method = RequestMethod.GET)
    Blog blogDetail(@PathVariable Long id) {
        return blogService.get(id);
    }

    @ApiOperation("create blog")
    @RequestMapping(method = RequestMethod.POST)
    void create(String title, String content) {
        blogService.createBlog(title, content);
    }
}
