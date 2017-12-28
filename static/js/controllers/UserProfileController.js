angular.module('BoilerAdmin').controller('UserProfileController', function($rootScope, $scope, $http, $timeout, $state, $window, moment) {
    profile = this;

    $rootScope.$on('$viewContentLoading',
        function(event, viewConfig){
            // Access to all the view config properties.
            // and one special property 'targetView'
            // viewConfig.targetView
        });

    $scope.$on('$viewContentLoaded', function() {   
        App.initAjax(); // initialize core components
        Layout.setAngularJsSidebarMenuActiveLink('set', $('#sidebar_menu_link_profile'), $state); // set profile link active in sidebar menu

        if (!$rootScope.currentUser) {
            $rootScope.getCurrentUser(refreshData);
        } else {
            refreshData();
        }
    });

    var refreshData = function () {
        profile.aUsername = $rootScope.currentUser.Username;
        profile.aRole = $rootScope.currentUser.Role.RoleId;
        profile.aOrg = $rootScope.currentUser.Organization ? $rootScope.currentUser.Organization.Uid : "";

        profile.aName = $rootScope.currentUser.Name;
        profile.aMobile = $rootScope.currentUser.MobileNumber;
        profile.aEmail = $rootScope.currentUser.Email;

        profile.loginName = $rootScope.currentUser.Username;
        if ($rootScope.currentUser.Status == 0) {
            profile.loginName += "(未激活)";
        }

        profile.picture = $rootScope.settings.layoutPath + "/img/" + "avatar0.png";
        if ($rootScope.currentUser.Picture.indexOf("avatar") > -1) {
            profile.picture = $rootScope.settings.layoutPath + "/img/" + $rootScope.currentUser.Picture;
        } else {
            profile.picture = $rootScope.currentUser.Picture;
        }

        profile.roleName = $rootScope.currentUser.Role.Name + (
                $rootScope.isOrgs() ? ("|" + $rootScope.currentUser.Organization.Type.Name) : "");

        profile.hasWeixin = false;
        profile.weixinName = "未绑定";

        angular.forEach($rootScope.currentUser.Thirds, function (third, key) {
            if (third.Platform === "weixin") {
                profile.hasWeixin = true;
                profile.weixinName = third.Name;
                return;
            }
        });
    };

    profile.updateUser = function() {
        $http.post("/user_profile_update/", {
            fullname: profile.aName,
            mobile: profile.aMobile,
            email: profile.aEmail
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

    profile.updateUserPassword = function() {
        if (profile.password.length <= 0 || profile.password_new.length < 0 ||
            profile.password_new != profile.password_new_confirm) {
            return;
        }
        $http.post("/user_password_update/", {
            password: profile.password,
            password_new: profile.password_new,
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

    profile.resetUser = function() {
        refreshData();
    };

    profile.weixinClick = function () {
        if (profile.hasWeixin) {
            swal({
                title: "是否解除微信账号的绑定？",
                text: "",//"解绑后，该用户将不能用微信扫码登录",
                type: "question",
                showCancelButton: true,
                //confirmButtonClass: "btn-danger",
                confirmButtonColor: "#d33",
                cancelButtonText: "取消",
                confirmButtonText: "解除绑定",
                //closeOnConfirm: false
            }).then(function () {
                $http.post('/user_unbind_weixin/', {
                    uid: $rootScope.currentUser.Uid
                }).then(function (res) {
                    $rootScope.getCurrentUser(refreshData);
                    swal({
                        title: "您的平台用户已经与微信账号解绑",
                        type: "success"
                    })
                }, function (err) {
                    console.error('unbind failed!', err.status, err.data);
                });
            });
        } else {
            swal({
                title: "是否进行微信账号的绑定？",
                text: "", //"绑定后，该用户将不能用微信扫码登录",
                type: "question",
                showCancelButton: true,
                //confirmButtonClass: "btn-danger",
                //confirmButtonColor: "#d33",
                cancelButtonText: "取消",
                confirmButtonText: "绑定微信",
                //closeOnConfirm: false
            }).then(function () {
                $window.location = "/user_bind_weixin/";
                // $http.get('/user_bind_weixin/').then(function (res) {
                //     // swal({
                //     //     title: "您的平台用户已经与微信账号解绑",
                //     //     type: "success"
                //     // })
                //     $rootScope.getCurrentUser(refreshData);
                // }, function (err) {
                //     alert('unbind failed! (' + err.status + ')' + err.data);
                // });
            });
        }
    };

    profile.bindWeixin = function () {

    };

    // set sidebar closed and body solid layout mode
    $rootScope.settings.layout.pageBodySolid = true;
    $rootScope.settings.layout.pageSidebarClosed = false;
}); 


var profile;