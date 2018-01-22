angular.module('BoilerAdmin').controller('UserProfileController', function($rootScope, $scope, $http, $timeout, $state, $window, Upload, moment) {
    bProfile = this;

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
            $rootScope.getCurrentUser();
        } else {
            refreshData();
        }
    });

    $rootScope.$watch("currentUser", function () {
        refreshData();
    });

    var refreshData = function () {
        bProfile.aUsername = $rootScope.currentUser.Username;
        bProfile.aRole = $rootScope.currentUser.Role.RoleId;
        bProfile.aOrg = $rootScope.currentUser.Organization ? $rootScope.currentUser.Organization.Uid : "";

        bProfile.aName = $rootScope.currentUser.Name;
        bProfile.aMobile = $rootScope.currentUser.MobileNumber;
        bProfile.aEmail = $rootScope.currentUser.Email;

        bProfile.loginName = $rootScope.currentUser.Username;
        if ($rootScope.currentUser.Status == 0) {
            bProfile.loginName += "(未激活)";
        }

        bProfile.picture = $rootScope.settings.layoutPath + "/img/" + "avatar0.png";
        if ($rootScope.currentUser.Picture.indexOf("avatar") > -1) {
            bProfile.picture = $rootScope.settings.layoutPath + "/img/" + $rootScope.currentUser.Picture;
        } else {
            bProfile.picture = $rootScope.currentUser.Picture;
        }

        bProfile.roleName = $rootScope.currentUser.Role.Name + (
                $rootScope.isOrgs() ? ("|" + $rootScope.currentUser.Organization.Type.Name) : "");

        bProfile.hasWeixin = false;
        bProfile.weixinName = "未绑定";

        angular.forEach($rootScope.currentUser.Thirds, function (third, key) {
            if (third.Platform === "weixin") {
                bProfile.hasWeixin = true;
                bProfile.weixinName = third.Name;
                return;
            }
        });
    };

    bProfile.avatarChanged = function () {
        console.warn("Init avatarChanged");
        console.warn(bProfile.avatar);
    };

    bProfile.removeAvatar = function () {
        bProfile.avatar = undefined;
    };

    bProfile.updateUser = function() {
        $http.post("/user_profile_update/", {
            fullname: bProfile.aName,
            mobile: bProfile.aMobile,
            email: bProfile.aEmail
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

    // upload on file select or drop
    bProfile.uploadAvatar = function () {
        if (!bProfile.avatar) {
            return;
        }

        Upload.upload({
            url: '/user_image_upload/',
            data: {
                uid: $rootScope.currentUser.Uid,
                file: bProfile.avatar,
            }
        }).then(function (resp) {
            console.log('Success ' + resp.config.data.file.name + 'uploaded. Response: ' + resp.data);
            swal({
                title: "头像上传成功",
                type: "success"
            }).then(function () {
                $rootScope.getCurrentUser();
            });
        }, function (err) {
            console.log('Error status: ' + err.status);
            swal({
                title: "头像上传失败",
                text: err.data,
                type: "error"
            });
        }, function (evt) {
            var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
            console.log('progress: ' + progressPercentage + '% ' + evt.config.data.file.name);
        });
    };

    bProfile.updateUserPassword = function() {
        if (bProfile.password.length <= 0 || bProfile.password_new.length < 0 ||
            bProfile.password_new != bProfile.password_new_confirm) {
            return;
        }
        $http.post("/user_password_update/", {
            password: bProfile.password,
            password_new: bProfile.password_new,
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

    bProfile.resetUser = function() {
        refreshData();
    };

    bProfile.weixinClick = function () {
        if (bProfile.hasWeixin) {
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

    bProfile.bindWeixin = function () {

    };

    // set sidebar closed and body solid layout mode
    $rootScope.settings.layout.pageBodySolid = true;
    $rootScope.settings.layout.pageSidebarClosed = false;
}); 


var bProfile;