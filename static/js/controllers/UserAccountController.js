/* Setup general page controller */
angular.module('BoilerAdmin').controller('UserAccountController', function($rootScope, $scope, $http, $uibModal, settings, moment, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    bAccount = this;
    bAccount.isDone = false;

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        $http.get('/user_roles/')
            .then(function (res) {
                bAccount.roles = res.data;
                bAccount.init();
            }, function (err) {
                console.error("Get Roles List Err: ", err);
            });

        // set default layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    bAccount.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers')
        .withOption('rowCallback', rowCallback);

    bAccount.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6).notSortable()
    ];

    bAccount.editing = false;

    bAccount.status = [
        //{id: -1, name: "新用户", hidden: true},
        {id: 0, name: "未激活", hidden: true},
        {id: 1, name: "通常"},
        {id: 2, name: "禁用"}
    ];

    bAccount.refreshDataTables = function () {
        $http.get('/user_list/')
            .then(function (res) {
                var data = res.data;

                var num = 0;
                angular.forEach(data, function (d, key) {
                    d.num = ++num;
                    d.stat = bAccount.status[d.Status];
                });

                bAccount.datasource = data;

                bAccount.isDone = true;
                setTimeout(function () {
                    App.stopPageLoading();
                }, 1500);
            });
    };

    var someClickHandler = function(info) {
        bAccount.message = info.Uid + ' - ' + info.Name;
    };

    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        $('td', nRow).unbind('click');
        $('td', nRow).bind('click', function() {
            bAccount.editing = false;
            bAccount.currentData = bAccount.datasource[aData[0] - 1];
            bAccount.currentData.aName = bAccount.currentData.Name;
            bAccount.currentData.aPassword = "";
            bAccount.currentData.resetPassowrd = false;
            bAccount.currentData.aRole = bAccount.currentData.Role.RoleId;
            bAccount.currentData.aStat = bAccount.currentData.Status;
            bAccount.currentData.aOrg = bAccount.currentData.Organization ? bAccount.currentData.Organization.Uid : "";
            // $scope.$apply(function() {
            //     someClickHandler(bAccount.currentData);
            // });
            var modalInstance = $uibModal.open({
                templateUrl: 'myModalContent.html',
                controller: 'userCtrl',
                size: "",
                windowClass: 'zindex',
                resolve: {
                    currentData: function () {
                        return bAccount.currentData;
                    },
                    status:function () {
                        return bAccount.status;
                    },
                    aRoles:function () {
                        return bAccount.aRoles;
                    }

                }
            });

            modalInstance.result.then(function (selectedItem) {
                $scope.selected = selectedItem;
                bAccount.refreshDataTables();
            }, function () {

            });

        });
        return nRow;
    }

    bAccount.init = function () {
        bAccount.aRoles = [];
        if ($rootScope.currentUser.Role.RoleId < 10) {
            angular.forEach(bAccount.roles, function (d, key) {
                if (d.RoleId > $rootScope.currentUser.Role.RoleId) {
                    bAccount.aRoles.push({ id: d.RoleId, name: d.Name });
                }
            });
        }

        if (Math.floor($rootScope.currentUser.Role.RoleId / 10) === 1) {
            angular.forEach(bAccount.roles, function (d, key) {
                if (d.RoleId > $rootScope.currentUser.Role.RoleId && d.RoleId < 20) {
                    bAccount.aRoles.push({ id: d.RoleId, name: d.Name });
                }
            });
        }

        if ($rootScope.currentUser.Role.RoleId >= 20) {
            angular.forEach(bAccount.roles, function (d, key) {
                if (d.RoleId >= $rootScope.currentUser.Role.RoleId) {
                    bAccount.aRoles.push({ id: d.RoleId, name: d.Name });
                }
            });
        }
    };

    bAccount.aStatus = [];

    bAccount.editRow = function() {
        if (!bAccount.currentData) {
            return;
        }

        bAccount.editing = true;
    };

    bAccount.new = function () {
        bAccount.currentData = null;
        bAccount.open();
    };

    bAccount.isOrgs = function () {
        return bAccount.currentData && bAccount.currentData.aRole > 1;
    };

    bAccount.activeRow = function() {
        var aData = bAccount.currentData;

        $http.post("/user_active/", {
            uid: aData.Uid
        }).then(function (res) {
            swal({
                title: "用户" + aData.Username + "激活成功",
                type: "success"
            }).then(function () {
                bAccount.refreshDataTables();
            });
        }, function (err) {
            swal({
                title: "用户" + aData.Username + "激活失败",
                text: err.data,
                type: "error"
            });
        });

        //oTable.fnDraw();
    };

    bAccount.resetPassword = function () {
        if (bAccount.currentData) {
            bAccount.currentData.resetPassword = true;
        }
    };

    bAccount.deleteRow = function() {
        var aData = bAccount.currentData;
        swal({
            title: "确认删除用户" + aData.Username + "？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/user_delete/", {
                uid: aData.Uid
            }).then(function (res) {
                bAccount.refreshDataTables();
                swal({
                    title: "用户" + aData.Username + "删除成功",
                    type: "success"
                }).then(function () {
                    // var idx = bAccount.datasource.indexOf(aData);
                    // if (idx > -1) {
                    //     bAccount.datasource.splice(idx, 1);
                    // }

                });
            }, function (err) {
                swal({
                    title: "删除用户失败",
                    text: err.data,
                    type: "error"
                });
            });
        });


        //oTable.fnDraw();
    };

    bAccount.resetRow = function () {
        bAccount.editing = false;
        bAccount.currentData.aName = bAccount.currentData.Name;
        bAccount.currentData.aPassword = "";
        bAccount.currentData.resetPassowrd = false;
        bAccount.currentData.aRole = bAccount.currentData.Role.RoleId;
        bAccount.currentData.aStat = bAccount.currentData.Status;
        bAccount.currentData.aOrg = bAccount.currentData.Organization ? bAccount.currentData.Organization.Uid : "";
    };

    bAccount.saveRow = function() {
        var aData = bAccount.currentData;

        var org = '';
        if (bAccount.isOrgs()) {
            org = aData.aOrg;
        }

        var data = {
            uid: aData.Uid,
            //username: username,
            fullname: aData.aName,
            role: aData.aRole,
            stat: aData.aStat,
            org: org
        };

        if (aData.aPassword && aData.aPassword.length > 0) {
            data.password_new = aData.aPassword;
        }

        $http.post("/user_update/", data)
            .then(function (res) {
                bAccount.refreshDataTables();
                swal({
                    title: "用户" + aData.Username + "信息修改成功",
                    type: "success"
                }).then(function () {

                });
        }, function (err) {
            swal({
                title: "修改用户信息失败",
                text: err.data,
                type: "error"
            });
        });
    };

    bAccount.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/user_account_info.html',
            controller: 'ModalAccountCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                currentData: function () {
                    return bAccount.currentData;
                },
                roles: function () {
                    return bAccount.aRoles;
                }
            }
        });
    };
});

var bAccount;

angular.module('BoilerAdmin').controller('ModalAccountCtrl', function ($uibModalInstance, $rootScope, $http, $log, currentData, roles) {
    var $modal = this;

    $modal.isValid = false;
    $modal.data = {};
    $modal.roles = roles;

    console.warn("init ModalAccountCtrl with roles:", roles);
    if ($modal.roles.length === 1 && $modal.roles[0]) {
        $modal.data.role = $modal.roles[0].id;
    }
    if ($rootScope.currentUser.Role.RoleId > 1) {
        $modal.data.org = $rootScope.currentUser.Organization.Uid;
    }


    $modal.dataChanged = function () {
        if ($modal.data.username.length < 6 || $modal.data.username.length > 16 ||
            $modal.data.password.length < 6 || $modal.data.username.length > 16 ||
            !$modal.data.role ||
            $modal.data.role <= $rootScope.currentUser.Role.RoleId ||
            ($modal.data.role > 1 && $modal.data.org.length <= 0)) {
            $modal.isValid = false;

            return;
        }

        $modal.isValid = true;
    };

    $modal.commit = function () {
        if (!$modal.isValid) {
            return
        }

        Ladda.create(document.getElementById('boiler_ok')).start();
        $modal.data.uid = "";
        $http.post("/user_update/", $modal.data)
            .then(function (res) {
                bAccount.refreshDataTables();
                swal({
                    title: "用户" + $modal.data.username + "添加成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.dismiss('cancel');
                });
            }, function (err) {
                swal({
                    title: "添加用户失败",
                    text: err.data,
                    type: "error"
                });
            });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.delete = function () {

    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };
});

// Please note that the close and dismiss bindings are from $uibModalInstance.
angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: '/directives/modal/terminal_config.html',
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


angular.module('BoilerAdmin').controller('userCtrl', function($scope,$rootScope, $uibModalInstance,$http, currentData,status,aRoles) {
    $scope.currentData= currentData;
    $scope.editing = false;
    $scope.currentUser = $rootScope.currentUser;
    $scope.status = status;
    $scope.aRoles = aRoles;

    //在这里处理要进行的操作
    $scope.saveRow = function() {
        var aData = $scope.currentData;
        var isOrgs = function () {
            return $scope.currentData && Math.floor($scope.currentData.aRole / 10) === 1;
        };
        var org = '';
        if (isOrgs()) {
            org = aData.aOrg;
        }

        var data = {
            uid: aData.Uid,
            //username: username,
            fullname: aData.aName,
            role: aData.aRole,
            stat: aData.aStat,
            org: org
        };

        if (aData.aPassword && aData.aPassword.length > 0) {
            data.password_new = aData.aPassword;
        }

        $http.post("/user_update/", data)
            .then(function (res) {
                bAccount.refreshDataTables();
                swal({
                    title: "用户" + aData.Username + "信息修改成功",
                    type: "success"
                }).then(function () {

                });
            }, function (err) {
                swal({
                    title: "修改用户信息失败",
                    text: err.data,
                    type: "error"
                });
            });
        $uibModalInstance.close();
    };
    $scope.resetRow = function() {
        $scope.editing = false;
        $scope.currentData.aName = $scope.currentData.Name;
        $scope.currentData.aPassword = "";
        $scope.currentData.resetPassowrd = false;
        $scope.currentData.aRole = $scope.currentData.Role.RoleId;
        $scope.currentData.aStat = $scope.currentData.Status;
        $scope.currentData.aOrg = $scope.currentData.Organization ? $scope.currentData.Organization.Uid : "";

    };
    $scope.activeRow = function(){
        var aData = $scope.currentData;

        $http.post("/user_active/", {
            uid: aData.Uid
        }).then(function (res) {
            swal({
                title: "用户" + aData.Username + "激活成功",
                type: "success"
            }).then(function () {
//	                bAccount.refreshDataTables();
            });
        }, function (err) {
            swal({
                title: "用户" + aData.Username + "激活失败",
                text: err.data,
                type: "error"
            });
        });
    };

    $scope.editRow = function() {

        if (!$scope.currentData) {
            return;
        }

        $scope.editing = true;
    };
    $scope.deleteRow = function(){
        var aData = $scope.currentData;
        swal({
            title: "确认删除用户" + aData.Username + "？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/user_delete/", {
                uid: aData.Uid
            }).then(function (res) {
                swal({
                    title: "用户" + aData.Username + "删除成功",
                    type: "success"
                }).then(function () {
                    // var idx = bAccount.datasource.indexOf(aData);
                    // if (idx > -1) {
                    //     bAccount.datasource.splice(idx, 1);
                    // }
                    bAccount.refreshDataTables();
                });
            }, function (err) {
                swal({
                    title: "删除用户失败",
                    text: err.data,
                    type: "error"
                });
            });
            $uibModalInstance.close();
        });

    }
    $scope.close = function(){
        $uibModalInstance.dismiss();
    }
    $scope.resetPassword=function(){
        if ($scope.currentData) {
            $scope.currentData.resetPassword = true;
        }
    };
});

