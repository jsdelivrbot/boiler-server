/***
 Metronic AngularJS App Main Script
 ***/

/* Boiler App */
var boilerAdmin = angular.module("BoilerAdmin", [
    "ui.router",
    "ui.bootstrap",
    "oc.lazyLoad",
    "ngSanitize",
    "ngCookies",
    "customFilter",
    "angularMoment",
    "ui.select",
    "datatables",
    "frapontillo.bootstrap-switch"
]);

/* Configure ocLazyLoader(refer: https://github.com/ocombe/ocLazyLoad) */
boilerAdmin.config(['$ocLazyLoadProvider', function($ocLazyLoadProvider) {
    $ocLazyLoadProvider.config({
        // global configs go here
        debug: true,
        events: true
    });
}]);

//AngularJS v1.3.x workaround for old style controller declarition in HTML
boilerAdmin.config(['$controllerProvider', function($controllerProvider) {
    // this option might be handy for migrating old apps, but please don't use it
    // in new ones!
    $controllerProvider.allowGlobals();
}]);

/********************************************
 END: BREAKING CHANGE in AngularJS v1.3.x:
 *********************************************/

var allcookies = document.cookie;
function getCookie(cookie_name) {
    var allcookies = document.cookie;
    var cookie_pos = allcookies.indexOf(cookie_name);

    if (cookie_pos != -1) {
        cookie_pos += cookie_name.length + 1;
        var cookie_end = allcookies.indexOf(";", cookie_pos);

        if (cookie_end == -1) {
            cookie_end = allcookies.length;
        }

        var value = unescape(allcookies.substring(cookie_pos, cookie_end));
    }
    return value;
}

/* Setup global settings */
boilerAdmin.factory('settings', ['$rootScope', function($rootScope) {
    // supported languages
    var settings = {
        layout: {
            pageSidebarClosed: false, // sidebar menu state
            pageContentWhite: true, // set page content layout
            pageBodySolid: false, // solid body color state
            pageAutoScrollOnLoad: 1000 // auto scroll to top on page load
        },
        assetsPath: '../assets',
        globalPath: '../assets/global',
        layoutPath: '../assets/layouts/boiler',
    };

    $rootScope.settings = settings;

    return settings;
}]);

const IP_JSON_URL = 'http://ipv4.myexternalip.com/json';

var clientkey = "?clientkey=" + getCookie("request_session");

/* Setup App Main Controller */
boilerAdmin.controller('AppController', ['$scope', '$rootScope', '$http', '$log', function($scope, $rootScope, $http, $log) {
    $scope.$on('$viewContentLoaded', function() {
        //App.initComponents(); // init core components
        //Layout.init(); //  Init entire layout(header, footer, sidebar, etc) on page load if the partials included in server side instead of loading with ng-include directive

    });
}]);

/***
 Layout Partials.
 By default the partials are loaded through AngularJS ng-include directive. In case they loaded in server side(e.g: PHP include function) then below partial
 initialization can be disabled and Layout.init() should be called on page load complete as explained above.
 ***/

var getZh = function (str) {
    var r = /\\u([\d\w]{4})/gi;
    var res = str.replace(r, function (match, grp) {
        return String.fromCharCode(parseInt(grp, 16));
    });

    return decodeURIComponent(res);
};


angular.module('BoilerAdmin').controller('ModalLoginCtrl', function ($uibModalInstance, $rootScope, $scope, $http) {
    var $modal = this;
    $modal.editing = false;
    $modal.roleId = 20;

    $modal.organizations = [];

    setTimeout(function () {
        $http.get('/organization_list/?scope=register')
            .then(function (res) {
                for (var i = 0; i < res.data.length; i++) {
                    var d = res.data[i];
                    d.name = d.Name;
                    d.type = d.Type.Name;
                    $modal.organizations.push(d);
                }
            }, function (err) {
                console.log("Get Register OrgList Err: " + err);
            });

        $http.get(IP_JSON_URL).then(function (result) {
            console.log("IP:" + result.data.ip);
            $modal.ip = result.data.ip;
        }, function (e) {
            console.error("Get IP Error:", e);
        });
    }, 0);

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };

    $modal.gotoSignup = function () {
        $uibModalInstance.dismiss('cancel');

        header.openSignup();
    };

    $modal.gotoLogin = function () {
        $uibModalInstance.dismiss('cancel');

        header.openLogin();
    };

    $modal.signup = function() {
        $http.post('/user_register_bind_third/', {
            username: $modal.username,
            password: $modal.password,
            mobile: $modal.mobile,
            //email: $scope.email,
            //name: $scope.fullname,
            role: $modal.roleId,
            // org: $modal.aOrg.Uid,
            ip: $modal.ip
        }).then(function (res) {
            $('#signup-form').modal('hide');
            swal({
                title: "注册成功",
                text: "您的平台账号 " + $modal.username + " 已经绑定微信，之后您可以通过用户名和密码进行登录，或者使用微信扫码直接登录平台，\n现在将转到该用户登录",
                type: "success",
                //showCancelButton: true,
                //confirmButtonClass: "btn-danger",
                confirmButtonText: "好的"
                //closeOnConfirm: false
            }).then(function () {
                $uibModalInstance.close('success');
                $rootScope.currentUser.Status = 1;
                header.refresh();
            });
        }, function (err) {
            var message = err.data;
            swal({
                title: "注册失败",
                text: message + "\n请返回重新填写",
                type: "warning",
                //showCancelButton: true,
                //confirmButtonClass: "btn-danger",
                confirmButtonText: "确定 "
            });
            this.remark = err.data;
        });

        //alert("register!\n" + data);
    };

    $modal.login = function() {
        var ip = "";
        $http.get(IP_JSON_URL).then(function(result) {
            console.log("ip" + result.data.ip);
            ip = result.data.ip;
        }, function(e) {
            console.error("Get IP Error:", e);
        }).then(function () {
            $http.post('/user_login_bind_third/', {
                username: $modal.username,
                password: $modal.password,
                ip: ip
            }).then(function (res) {
                $rootScope.getCurrentUser(function () {
                    header.refresh();
                    var text = "欢迎回来，" + $rootScope.currentUser.Role.Name + " " + $rootScope.currentUser.Username + "。";
                    text += "\n您的微信账号已经绑定成功，之后您可以通过微信扫码直接登录平台。";
                    swal({
                        title: "登录成功",
                        text: text,
                        type: "success",
                        //showCancelButton: true,
                        //confirmButtonClass: "btn-danger",
                        //cancelButtonText: "留在首页",
                        confirmButtonText: "好的",
                        //closeOnConfirm: false
                    }).then(function () {
                        $uibModalInstance.close('success');
                        header.refresh();
                    }, function (dismiss) {
                        //alert("login and refresh");
                        //refresh();
                    });
                });
            }, function (err) {
                swal({
                    title: "登录失败",
                    text: err.data,
                    type: "error",
                    //showCancelButton: true,
                    //confirmButtonClass: "btn-danger",
                    confirmButtonText: "确定"
                });
            });
        });
    };


});

// Please note that the close and dismiss bindings are from $uibModalInstance.

angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: 'myModalContent.html' + clientkey,
    bindings: {
        resolve: '<',
        close: '&',
        dismiss: '&'
    },
    controller: function () {
        var $ctrl = this;

        $ctrl.$onInit = function () {
            // $ctrl.items = $ctrl.resolve.items;
            // $ctrl.selected = {
            //     item: $ctrl.items[0]
            // };
        };

        $ctrl.ok = function () {
            // $ctrl.close({$value: $ctrl.selected.item});
        };

        $ctrl.cancel = function () {
            $ctrl.dismiss({$value: 'cancel'});
        };
    }
});

/* Setup Layout Part - Sidebar */
boilerAdmin.controller('SidebarController', ['$state', '$scope', function($state, $scope) {
    $scope.$on('$includeContentLoaded', function() {
        Layout.initSidebar($state); // init sidebar
    });
}]);

/* Setup Layout Part - Footer */
boilerAdmin.controller('FooterController', ['$scope', function($scope) {
    $scope.$on('$includeContentLoaded', function() {
        Layout.initFooter(); // init footer
    });
}]);

var thumb;

/* Setup Rounting For All Pages */
boilerAdmin.config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider) {
    // Redirect any unmatched url
    $urlRouterProvider.otherwise("/monitor/thumb");
    // console.warn("$stateProvider init!");
    /*============= MONITOR BEGIN   =============*/
    $stateProvider
        .state('monitor', {
            url: "/monitor",
            templateUrl: "views/monitor/main.html" + clientkey,
            data: {pageTitle: '主监控台'},
            controller: "DashboardController",
            controllerAs: "dashboard",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/global/plugins/jquery.sparkline.min.js' + clientkey,
                            '../assets/global/plugins/counterup/jquery.waypoints.min.js' + clientkey,
                            '../assets/global/plugins/counterup/jquery.counterup.min.js' + clientkey,

                            'js/controllers/DashboardController.js' + clientkey,

                            'js/directives/components/filter_monitor.js' + clientkey,
                            'js/directives/components/paging-bar.js' + clientkey,
                            'js/directives/location-select.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state('monitor.dashboard', {
            url: "/dashboard",
            templateUrl: "views/monitor/dashboard.html" + clientkey,
            data: {pageTitle: '平台总览'}
        })
        .state('monitor.thumb', {
            url: "/thumb",
            templateUrl: "views/monitor/thumb.html" + clientkey,
            data: {pageTitle: '设备图文'}
        })
        .state('monitor.list', {
            url: "/list",
            templateUrl: "views/monitor/list.html" + clientkey,
            data: {pageTitle: '设备列表'}
        })
        .state('monitor.map', {
            url: "/map",
            templateUrl: "views/monitor/map.html" + clientkey,
            data: {pageTitle: '设备地图'}
        });
    /*============= MONITOR END     =============*/

    /*============= RUNTIME BEGIN   =============*/
    $stateProvider
        .state('runtime', {
            url: "/runtime?:boiler:from",
            templateUrl: "views/runtime/main.html" + clientkey,
            data: {pageTitle: '主监控台 - 实时监控'},
            controller: "BoilerRuntimeController",
            controllerAs: "runtime",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'MainContent',
                        //insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            //'js/controllers/BoilerModuleController.js' + clientkey,

                            'js/directives/boiler_module.js' + clientkey,
                            'js/directives/chart_steam.js' + clientkey,
                            'js/directives/chart_smoke-components.js' + clientkey,
                            'js/directives/chart_temperature.js' + clientkey,
                            'js/directives/chart_excess-air.js' + clientkey,
                            'js/directives/chart_heat.js' + clientkey,
                            'js/directives/chart_heat_month.js' + clientkey,

                            'js/directives/chart_dynamic.js' + clientkey,

                            'js/directives/chart_status1.js' + clientkey,
                            'js/directives/chart_status2.js' + clientkey,
                            'js/directives/chart_status3.js' + clientkey,

                            'js/controllers/BoilerRuntimeController.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state("runtime.dashboard", {
            url: "/dashboard",
            templateUrl: "views/runtime/dashboard.html" + clientkey,
            data: {pageTitle: '炉型详图'}
        })

        .state("runtime.stats", {
            url: "/stats",
            templateUrl: "views/runtime/stats.html" + clientkey,
            data: {pageTitle: '运行参数'}
        })

        .state("runtime.history", {
            url: "/history",
            templateUrl: "views/runtime/history.html" + clientkey,
            controller: "BoilerHistoryController",
            data: {pageTitle: '历史数据'},
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerChartCore',
                        //insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            'bower_components/alasql/dist/alasql.min.js' + clientkey,
                            'bower_components/js-xlsx/dist/xlsx.full.min.js' + clientkey,

                            'js/controllers/BoilerHistoryController.js' + clientkey,
                            'js/directives/datatable_history.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state("runtime.info", {
            url: "/info",
            templateUrl: "views/runtime/info.html" + clientkey,
            controller: "BoilerInfoController",
            controllerAs: "info",
            data: {pageTitle: '设备信息'},
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerChartCore',
                        //insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,

                            'js/controllers/BoilerInfoController.js' + clientkey,
                            'js/directives/table_boiler-info.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state("runtime.alarm", {
            url: "/alarm",
            templateUrl: "views/runtime/alarm.html" + clientkey,
            controller: "AlarmController",
            controllerAs: "alarm",
            data: {pageTitle: '锅炉告警'},
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerChartCore',
                        //insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,

                            'js/controllers/AlarmController.js' + clientkey,
                            'js/directives/datatable_alarm.js' + clientkey,
                            'js/directives/chart_alarm.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state("runtime.maintain", {
            url: "/maintain",
            templateUrl: "views/runtime/maintain.html" + clientkey,
            data: {pageTitle: '维保记录'},
            controller: "BoilerMaintainController",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before',
                        files: [
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,

                            'js/controllers/BoilerMaintainController.js' + clientkey,
                            'js/directives/datatable_maintain.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state("runtime.developer", {
            url: "/developer",
            templateUrl: "views/runtime/developer.html" + clientkey,
            data: {pageTitle: '调试设置'},
            controller: "BoilerDeveloperController",
            controllerAs: "developer",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        files: [
                            'js/controllers/BoilerDeveloperController.js' + clientkey,
                        ]
                    });
                }]
            }
        });
    /*============= RUNTIME END     =============*/

    // Alarm
    $stateProvider
        .state('alarm', {
            url: "/alarm",
            templateUrl: "views/alarm.html" + clientkey,
            data: {pageTitle: '告警信息'},
            controller: "AlarmController",
            controllerAs: "alarm",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/pages/scripts/table-datatables-managed.js' + clientkey,

                            '../assets/boiler/global/plugins/moment.js' + clientkey,
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/AlarmController.js' + clientkey,
                            'js/directives/datatable_alarm.js' + clientkey,
                            'js/directives/chart_alarm.js' + clientkey
                        ]
                    });
                }]
            }
        })

        // 专家咨询
        .state('dialogue', {
            url: "/dialogue",
            templateUrl: "views/dialogue.html" + clientkey,
            data: {pageTitle: '专家咨询'},
            controller: "DialogueController",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/pages/scripts/table-datatables-managed.js' + clientkey,

                            '../assets/boiler/global/plugins/moment.js' + clientkey,
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/DialogueController.js' + clientkey,
                        ]
                    });
                }]
            }
        })

        // 维保记录
        .state('boiler-maintain', {
            url: "/boiler-maintain",
            templateUrl: "views/boiler-maintain.html" + clientkey,
            data: {pageTitle: '维保记录'},
            controller: "BoilerMaintainController",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/BoilerMaintainController.js' + clientkey,
                            'js/directives/datatable_maintain.js' + clientkey
                        ]
                    });
                }]
            }
        })

        // 告警设置
        .state('config-runtime-alarm', {
            url: "/config-runtime-alarm",
            templateUrl: "views/config-runtime-alarm.html" + clientkey,
            data: {pageTitle: '告警设置'},
            controller: "ConfigAlarmController",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/pages/scripts/table-datatables-managed.js' + clientkey,

                            '../assets/boiler/global/plugins/moment.js' + clientkey,
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/ConfigAlarmController.js' + clientkey,
                        ]
                    });
                }]
            }
        });

    /*============= PARAMETER BEGIN   =============*/
    $stateProvider
        .state('parameter', {
            url: "/parameter",
            templateUrl: "views/parameter/main.html" + clientkey,
            data: {pageTitle: '运行时参数'},
            controller: "ParameterController",
            controllerAs: "param",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'MainContent',
                        //insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/pages/scripts/table-datatables-managed.js' + clientkey,

                            '../assets/boiler/global/plugins/moment.js' + clientkey,
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/ParameterController.js' + clientkey,
                        ]
                    });
                }]
            }
        })

        .state("parameter.dashboard", {
            url: "/dashboard",
            templateUrl: "views/parameter/dashboard.html" + clientkey,
            data: {pageTitle: '参数列表'}
        })
    ;
    /*============= PARAMETER END     =============*/


    /*============ ORGANIZATION BEGIN ============*/
    $stateProvider
        .state('organization', {
            url: "/organization?:tid",
            templateUrl: "views/organization/main.html" + clientkey,
            data: {pageTitle: '企业信息总览'},
            controller: "OrganizationController",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/global/plugins/moment.js' + clientkey,
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/OrganizationController.js' + clientkey,
                        ]
                    });
                }
                ]
            }
        })

        .state("organization.overview", {
            url: "/overview",
            templateUrl: "views/organization/dashboard.html" + clientkey,
            data: {pageTitle: '企业信息总览'}
        })

        .state("organization.enterprise", {
            url: "/enterprise",
            templateUrl: "views/organization/main.html" + clientkey,
            data: {pageTitle: '用能企业'}
        })

        .state("organization.factory", {
            url: "/factory",
            templateUrl: "views/organization/main.html" + clientkey,
            data: {pageTitle: '锅炉制造厂'}
        })

        .state("organization.installer", {
            url: "/installer",
            templateUrl: "views/organization/main.html" + clientkey,
            data: {pageTitle: '安装企业'}
        })

        .state("organization.government", {
            url: "/government",
            templateUrl: "views/organization/main.html" + clientkey,
            data: {pageTitle: '政府机关'}
        })

        .state("organization.supervisor", {
            url: "/supervisor",
            templateUrl: "views/organization/main.html" + clientkey,
            data: {pageTitle: '监管部门'}
        })

        /*============ ORGANIZATION END ============*/

        // Account
        .state('user-account', {
            url: "/user-account",
            templateUrl: "views/user-account.html" + clientkey,
            data: {pageTitle: '账号管理'},
            controller: "UserAccountController",
            controllerAs: "account",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/pages/scripts/table-datatables-managed.js' + clientkey,

                            '../assets/boiler/global/plugins/moment.js' + clientkey,
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/UserAccountController.js' + clientkey,
                        ]
                    });
                }]
            }
        })


        /*============= BOILER BEGIN =============*/
        .state('boiler', {
            url: "/boiler",
            templateUrl: "views/boiler/main.html" + clientkey,
            data: {pageTitle: '锅炉信息'},
            controller: "BoilerInfoController",
            controllerAs: "info",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerInfoList',
                        //insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,

                            'js/controllers/BoilerInfoController.js' + clientkey,
                            'js/directives/table_boiler-info.js' + clientkey
                        ]
                    });
                }]
            }
        })

        .state("boiler.dashboard", {
            url: "/dashboard",
            templateUrl: "views/boiler/dashboard.html" + clientkey,
            data: {pageTitle: '设备详细信息'},
            // resolve: {
            //     deps: ['$ocLazyLoad', function($ocLazyLoad) {
            //         return $ocLazyLoad.load({
            //             name: 'BoilerContent',
            //             files: [
            //                 'js/directives/table_boiler-info.js' + clientkey
            //             ]
            //         });
            //     }]
            // }
        })

        .state("boiler.info", {
            url: "/info?:boiler:from",
            templateUrl: "views/boiler/info.html" + clientkey,
            data: {pageTitle: '锅炉详情'}
        })
        /*============= BOILER END =============*/

        /*============= TERMINAL BEGIN =============*/
        .state('terminal', {
            url: "/terminal",
            templateUrl: "views/terminal/main.html" + clientkey,
            data: {pageTitle: '终端管理'},
            controller: "TerminalController",
            controllerAs: "terminal",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before a LINK element with this ID. Dynamic CSS files must be loaded between core and theme css files
                        files: [
                            '../assets/boiler/pages/scripts/table-datatables-managed.js' + clientkey,

                            '../assets/boiler/global/plugins/fullcalendar/fullcalendar.js' + clientkey,
                            '../assets/global/plugins/jquery-ui/jquery-ui.min.js' + clientkey,

                            'js/controllers/TerminalController.js' + clientkey,
                        ]
                    });
                }]
            }
        })

        .state('terminal.dashboard', {
            url: "/dashboard",
            templateUrl: "views/terminal/dashboard.html" + clientkey,
            data: {pageTitle: '终端列表'}
        })

        .state("terminal.message", {
            url: "/message?:terminal",
            templateUrl: "views/terminal/message.html" + clientkey,
            data: {pageTitle: '终端消息调试'}
        })

        .state("terminal.status", {
            url: "/status?:terminal",
            params:{"terminal":null},
            templateUrl: "views/terminal/status.html" + clientkey,
        })

        .state("terminal.configStatus", {
            url: "/configStatus",
            params:{"data":null},
            templateUrl: "views/terminal/configStatus.html" + clientkey,
        })

        /*============= TERMINAL END =============*/

        // User Profile
        .state("profile", {
            url: "/profile",
            templateUrl: "views/profile/main.html" + clientkey,
            data: {pageTitle: '用户信息'},
            controller: "UserProfileController",
            controllerAs: "profile",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before '#ng_load_plugins_before'
                        files: [
                            '../assets/global/plugins/bootstrap-fileinput/bootstrap-fileinput.css',
                            '../assets/pages/css/profile.css',

                            '../assets/global/plugins/jquery.sparkline.min.js' + clientkey,
                            '../assets/global/plugins/bootstrap-fileinput/bootstrap-fileinput.js' + clientkey,

                            '../assets/pages/scripts/profile.min.js' + clientkey,
                            'bower_components/ng-file-upload/ng-file-upload.min.js',

                            'js/controllers/UserProfileController.js' + clientkey
                        ]
                    });
                }]
            }
        })

        // User Profile Dashboard
        .state("profile.dashboard", {
            url: "/dashboard",
            templateUrl: "views/profile/dashboard.html" + clientkey,
            data: {pageTitle: '用户信息'}
        })

        // User Profile Account
        .state("profile.account", {
            url: "/account",
            templateUrl: "views/profile/account.html" + clientkey,
            data: {pageTitle: '账户设置'}
        })

        // User Profile Help
        .state("profile.help", {
            url: "/help",
            templateUrl: "views/profile/help.html" + clientkey,
            data: {pageTitle: '帮助'}
        })

        // Wiki
        .state("wiki", {
            url: "/wiki",
            templateUrl: "views/wiki/main.html" + clientkey,
            data: {pageTitle: '系统帮助'},
            controller: "WikiController",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before '#ng_load_plugins_before'
                        files: [
                            // '../assets/global/plugins/bootstrap-fileinput/bootstrap-fileinput.css',
                            // '../assets/pages/css/profile.css',

                            '../assets/global/plugins/jquery.sparkline.min.js' + clientkey,
                            '../assets/global/plugins/bootstrap-fileinput/bootstrap-fileinput.js' + clientkey,

                            '../assets/pages/scripts/profile.min.js' + clientkey,

                            'js/controllers/WikiController.js' + clientkey
                        ]
                    });
                }]
            }
        })

        // User Profile Dashboard
        .state("wiki.dashboard", {
            url: "/dashboard",
            templateUrl: "views/wiki/dashboard.html" + clientkey,
            data: {pageTitle: '帮助总览'}
        })
        .state("upload",{
            url:"/upload",
            templateUrl:"views/upload-file.html",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before '#ng_load_plugins_before'
                        files: [
                            'js/controllers/uploadFileController.js' + clientkey
                        ]
                    });
                }]
            }
        })
        .state("template",{
            url:"/template",
            templateUrl:"views/templates.html",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before '#ng_load_plugins_before'
                        files: [
                            'js/controllers/templateController.js' + clientkey
                        ]
                    });
                }]
            }
        })
        .state("wizard",{
            url:"/config-wizard",
            templateUrl:"views/config-wizard.html",
            resolve: {
                deps: ['$ocLazyLoad', function($ocLazyLoad) {
                    return $ocLazyLoad.load({
                        name: 'BoilerAdmin',
                        insertBefore: '#ng_load_plugins_before', // load the above css files before '#ng_load_plugins_before'
                        files: [
                            'js/controllers/configWizardController.js' + clientkey
                        ]
                    });
                }]
            }
        })

}]);

/* Init global settings and run the app */
boilerAdmin.run(["$rootScope", "$http", "$log", "$timeout", "settings", "$state", "$stateParams", "amMoment", "DTDefaultOptions",
    function($rootScope, $http, $log, $timeout, settings, $state, $stateParams, amMoment, DTDefaultOptions) {
        $rootScope.$state = $state;         // state to be accessed from view
        $rootScope.$stateParams = $stateParams;
        $rootScope.$settings = settings;    // state to be accessed from view
        amMoment.changeLocale('zh-cn');

        DTDefaultOptions.setLanguage({
            aria: {
                sortAscending: ": 激活正序排序",
                sortDescending: ": 激活反序排序"
            },
            emptyTable: "没有找到有效的记录",
            info: "第 _START_ 到 _END_ 共 _TOTAL_ 条记录",
            infoEmpty: "没有有效的记录",
            infoFiltered: "(筛选自总计 _MAX_ 条记录)",
            lengthMenu: "每页显示 _MENU_ 条",
            search: "搜索：",
            zeroRecords: "没有找到适合的记录",
            paginate: {
                previous:"上一页",
                next: "下一页",
                last: "末页",
                first: "首页"
            }
        });

        $rootScope.getBoilerList = function () {
            $http.get('/boiler_list/')
                .then(function (res) {
                    $log.log('Global BoilerList Get:', res);
                    $rootScope.boilers = res.data;
                    /*
                     $timeout(function () {
                     for (var i = 0; i < $rootScope.boilers.length; i++) {
                     var ab = $rootScope.boilers[i];
                     $rootScope.getBoilerCalculateParameter(ab);
                     }
                     }, 1000);
                     */
                    // $log.warn("Global Boilers:", $rootScope.boilers);

                    $timeout(function () {
                        $rootScope.getAlarmList();
                    }, 2000);
                }, function (err) {
                    $rootScope.boilers = [];
                });
        };

        $rootScope.getBoilerCalculateParameter = function (boiler) {
            $http.get('/boiler_calculate_parameter/?boiler=' + boiler.Uid)
                .then(function (res) {
                    boiler.Calculate = res.data;
                });
        };

        $rootScope.getAlarmList = function () {
            $http.get('/boiler_alarm_list/')
                .then(function (res) {
                    if (res.data) {
                        $rootScope.boilerAlarms = res.data;
                    } else {
                        $rootScope.boilerAlarms = [];
                    }
                });
        };

        $rootScope.getOrganizationList = function () {
            $http.get('/organization_list/')
                .then(function (res) {
                    // console.warn("Get OrganizationList: ", res);
                    var orgs = [];
                    for (var i in res.data) {
                        var d = res.data[i];
                        d.name = d.Name;
                        d.type = d.Type.Name;
                        d.typeId = d.Type.TypeId;
                        orgs.push(d);
                    }

                    $rootScope.organizations = orgs;
                }, function (err) {
                    //alert("Get OrganizationList Err: " + err);
                });

            $http.get('/organization_type_list/')
                .then(function (res) {
                    // console.warn("Get OrganizationList: ", res);
                    var types = [];
                    for (var i in res.data) {
                        var d = res.data[i];
                        var da = {};
                        da.id = d.TypeId;
                        da.name = d.Name;
                        types.push(da);
                    }

                    $rootScope.organizationTypes = types;
                }, function (err) {
                    //alert("Get OrganizationList Err: " + err);
                });
        };


        $http.get('/location_list/')
            .then(function (res) {
                $rootScope.locations = res.data;
            }, function (err) {
                //alert("Get LocationList Err: " + err);
            });


        $rootScope.getParameterList = function () {
            $http.get('/runtime_parameter_list/')
                .then(function (res) {
                    $rootScope.parameters = res.data;
                }, function (err) {
                    $rootScope.parameters = [];
                });
        };

        $http.get('/boiler_fuel_list/')
            .then(function (res) {
                $rootScope.fuels = res.data;
            });
        $http.get('/boiler_fuel_type_list/')
            .then(function (res) {
                $rootScope.fuelTypes = res.data;
            });
        $http.get('/boiler_form_list/')
            .then(function (res) {
                $rootScope.boilerForms = res.data;
            });
        $http.get('/boiler_medium_list/')
            .then(function (res) {
                $rootScope.boilerMediums = res.data;
            });

        $rootScope.getBoilerList();
        $rootScope.getParameterList();
        $rootScope.getOrganizationList();
    }]);