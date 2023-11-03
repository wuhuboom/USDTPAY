layui.config({  // common.js是配置layui扩展模块的目录，每个页面都需要引入
    version: true,   // 更新组件缓存，设为true不缓存，也可以设一个固定值
    base: getProjectUrl() + 'assets/module/'
}).extend({
    steps: 'steps/steps',
    notice: 'notice/notice',
    cascader: 'cascader/cascader',
    dropdown: 'dropdown/dropdown',
    fileChoose: 'fileChoose/fileChoose',
    Split: 'Split/Split',
    Cropper: 'Cropper/Cropper',
    tagsInput: 'tagsInput/tagsInput',
    citypicker: 'city-picker/city-picker',
    introJs: 'introJs/introJs',
    zTree: 'zTree/zTree',
    authtree: 'authtree',
    optimizeSelectOption:'optimizeSelectOption/optimizeSelectOption'
}).use(['layer', 'admin','table'], function () {
    var $ = layui.jquery;
    var layer = layui.layer;
    var admin = layui.admin;


    // /* table全局设置 */
    // var token = '这里你可以从缓存中获取token';
    // if (token && token.access_token) {
    //
    // }

    layui.table.set({
        parseData: function(res) {  // 利用parseData实现预处理
            // Sorry, your request is invalid
            if(res.code == -103 && res.msg.indexOf('Sorry') !== -1) { //token过期
                // setter.removeToken();
                layui.layer.msg('登录过期', {icon: 2, anim: 6, time: 1500}, function () {
                    // location.replace('components/template/login/login.html');
                });
            }
            return res;
        }
    });

});

/** 获取当前项目的根路径，通过获取layui.js全路径截取assets之前的地址 */
function getProjectUrl() {
    var layuiDir = layui.cache.dir;
    if (!layuiDir) {
        var js = document.scripts, last = js.length - 1, src;
        for (var i = last; i > 0; i--) {
            if (js[i].readyState === 'interactive') {
                src = js[i].src;
                break;
            }
        }
        var jsPath = src || js[last].src;
        layuiDir = jsPath.substring(0, jsPath.lastIndexOf('/') + 1);
    }
    return layuiDir.substring(0, layuiDir.indexOf('assets'));
}
