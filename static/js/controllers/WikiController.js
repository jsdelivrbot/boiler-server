angular.module('BoilerAdmin').controller('WikiController', function($rootScope, $scope, $http, $timeout, $state, moment) {
    wiki = this;

    $scope.$on('$viewContentLoaded', function() {   
        App.initAjax(); // initialize core components
        Layout.setAngularJsSidebarMenuActiveLink('set', $('#sidebar_menu_link_wiki'), $state); // set wiki link active in sidebar menu

        $http.get('/user_roles/')
            .then(function (res) {
                wiki.roles = res.data;
            }, function (err) {
                console.error("Get Roles List Err: ", err);
            });

    });

    var refreshData = function () {
        var roleNames = ["用户角色", "超级管理员", "管理员", "客服", "机构管理员", "机构用户", "普通用户", "未激活用户"];
        var canManageUserAccount = [""];
    };

    refreshData();

    wiki.loginName = $rootScope.currentUser.Username;
    if ($rootScope.currentUser.Status == 0) {
        wiki.loginName += "(未激活)";
    }
    wiki.roleName = $rootScope.currentUser.Role.Name + (
        $rootScope.isOrgs() ? ("|" + $rootScope.currentUser.Organization.Type.Name) : "");


    wiki.updateUser = function() {
        $http.post("/user_wiki_update/", {
            fullname: wiki.aName,
            mobile: wiki.aMobile,
            email: wiki.aEmail
        }).then(function (res) {
            swal({
                title: "您的个人信息修改成功",
                type: "success"
            }).then(function () {
                $rootScope.getCurrentUser();
            });
        }, function (err) {
            swal({
                title: "修改用户信息失败",
                text: err.data,
                type: "error"
            });
        });
    };

    wiki.updateUserPassword = function() {
        if (wiki.password.length <= 0 || wiki.password_new.length < 0 ||
            wiki.password_new != wiki.password_new_confirm) {
            return;
        }
        $http.post("/user_password_update/", {
            password: wiki.password,
            password_new: wiki.password_new,
        }).then(function (res) {
            swal({
                title: "密码修改成功",
                type: "success"
            }).then(function () {
                $rootScope.getCurrentUser();
            });
        }, function (err) {
            swal({
                title: "密码修改失败",
                text: err.data,
                type: "error"
            });
        });
    };

    wiki.resetUser = function() {
        refreshData();
    };

    // set sidebar closed and body solid layout mode
    $rootScope.settings.layout.pageBodySolid = true;
    $rootScope.settings.layout.pageSidebarClosed = false;
}); 


var wiki;