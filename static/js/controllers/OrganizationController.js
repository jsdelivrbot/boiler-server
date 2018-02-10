angular.module('BoilerAdmin').controller('OrganizationController', function($rootScope, $scope, $http, $uibModal, $log, $document, $location, $timeout, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    organization = this;
    organization.isDone = false;
    organization.tid = 0;

    $rootScope.$watch("organizations", function () {
        if (!$rootScope.organizations) {
            return;
        }

        organization.refreshDataTables();
    });

    var p = $location.search();
    if (!p['tid'] || p['tid'].length === 0) {
        p['tid'] = "";
    } else {
        organization.tid = parseInt(p['tid']);
    }

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    organization.titles = [
        '企业总表',
        '锅炉制造企业列表',
        '锅炉使用企业列表',
        '锅炉安装企业列表',
        '控制柜供应商列表',
        '锅炉维保企业列表',
        '能源管理企业列表',
        '政府监管部门列表'
    ];

    organization.refreshDataTables = function () {
        var orgs = [];
        for (var i in $rootScope.organizations) {
            var og = $rootScope.organizations[i];
            og.num = parseInt(i, 10) + 1;

            if (organization.tid > 0) {
                if (og.typeId == organization.tid) {
                    orgs.push(og);
                }
            } else {
                orgs.push(og);
            }
        }

        organization.datasource = orgs;
        organization.isDone = true;

        setTimeout(function () {
            App.stopPageLoading();
        }, 500);
    };

    organization.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
    // .withOption('rowCallback', rowCallback);

    organization.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4)
        // DTColumnDefBuilder.newColumnDef(5).notSortable()
    ];

    organization.new = function () {
        $log.info("organization.new");
        currentData = null;
        editing = true;
        organization.open('lg');
    };

    organization.edit = function (uid) {
        $log.info("organization.edit:", uid);
        for (var i = 0; i < organization.datasource.length; i++) {
            if (organization.datasource[i].Uid === uid) {
                currentData = organization.datasource[i];
                editing = true;
                $log.info("organization.edit GET:", currentData);
                organization.open('lg');
                break;
            }
        }
    };

    organization.view = function (uid) {
        $log.info("organization.edit:", uid);
        for (var i = 0; i < organization.datasource.length; i++) {
            if (organization.datasource[i].Uid === uid) {
                currentData = organization.datasource[i];
                editing = false;
                $log.info("organization.edit GET:", currentData);
                organization.open('lg');
                break;
            }
        }
    };

    organization.delete = function (o) {
        swal({
            title: "确认删除该企业？\n" + o.Name,
            text: "注意：删除后将无法恢复，且企业相关用户会一并删除。",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/organization_delete/", {
                uid: o.Uid
            }).then(function (res) {
                swal({
                    title: "企业删除成功",
                    text: "如需移除该企业相关锅炉，请在锅炉列表中进行删除，或联系管理员进行操作。",
                    type: "success"
                }).then(function () {
                    $rootScope.getOrganizationList();
                    organization.refreshDataTables();
                });
            }, function (err) {
                swal({
                    title: "企业删除失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

    organization.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/organization_detail.html',
            controller: 'ModalOrganizationCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            // resolve: {
            //     items: function () {
            //         return dialogue.items;
            //     }
            // }
        });

        modalInstance.result.then(function (selectedItem) {
            // dialogue.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };
});

var organization;
var currentData;
var editing;

angular.module('BoilerAdmin').controller('ModalOrganizationCtrl', function ($uibModalInstance, $rootScope, $http) {
    var $modal = this;
    $modal.editing = editing;
    $modal.org = currentData;
    $modal.title = "新增企业信息";
    $modal.typeId = -1;

    $modal.isSuper = false;
    $modal.supervisor = null;

    $modal.showBrand = false;
    $modal.brandName = "";

    $modal.aProvince = $rootScope.locations[0];
    $modal.aProvince.Name = "所在区域";
    $modal.location = $modal.aProvince;

    $rootScope.$watch("organizationTypes", function () {
        if (!$rootScope.organizationTypes) {
            return;
        }

        $modal.refreshOrganizationTypes();
    });

    $modal.refreshOrganizationTypes = function () {
        $modal.organizationTypes = $rootScope.organizationTypes;

        var def = {
            id: -1,
            name: "企业类型（请选择）"
        };

        if ($modal.organizationTypes.length <= 0 ||
            $modal.organizationTypes[0].id >= 0) {
            $modal.organizationTypes.unshift(def);
        }
    };

    $modal.changeProvince = function () {
        $modal.location = $modal.aProvince;
        // dashboard.filterBoilers();
    };

    $modal.changeCity = function () {
        $modal.location = $modal.aCity;
        // dashboard.filterBoilers();
    };

    $modal.changeRegion = function () {
        $modal.location = $modal.aRegion;
        // dashboard.filterBoilers();
    };

    var getLocation = function (locationId, locationList, locationScope) {
        for (var pi = 0; pi < locationList.length; pi++) {
            var local = locationList[pi];
            if (local.LocationId === Math.floor(locationId / 10000) ||
                local.LocationId === Math.floor(locationId / 100) ||
                local.LocationId === locationId) {
                switch (locationScope) {
                    case "province":
                        $modal.aProvince = local;
                        break;
                    case "city":
                        $modal.aCity = local;
                        break;
                    case "region":
                        $modal.aRegion = local;
                        break;
                }
                break;
            }
        }

        if (locationId < 100) {
            return;
        }

        switch (locationScope) {
            case "province":
                getLocation(locationId, $modal.aProvince.cities, "city");
                break;
            case "city":
                getLocation(locationId, $modal.aCity.regions, "region");
                break;
            case "region":
                break;
        }
    };

    if (currentData) {
        $modal.title = "企业信息";
        $modal.name = currentData.Name;
        $modal.typeId = currentData.Type.TypeId;
        $modal.address = currentData.Address.Address;
        $modal.location = currentData.Address.Location;

        $modal.showBrand = currentData.ShowBrand;
        $modal.brandName = currentData.BrandName;

        $modal.isSuper = currentData.IsSupervior;
        $modal.supervisor = currentData.SuperOrganization;

        var locationId = $modal.location.LocationId;
        getLocation(locationId, $rootScope.locations, "province");
    }

    /*
     MyUidObject

     Type			*OrganizationType	`orm:"rel(fk);index"`
     Address			*Address		`orm:"rel(fk);null"`
     Contact			*Contact		`orm:"rel(fk);null"`

     Users			[]*User			`orm:"reverse(many);null"`
     Boilers			[]*Boiler		`orm:"reverse(many);null"`
     */

    $modal.ok = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();
        var uid = currentData ? currentData.Uid : "";
        //alert("Ready to post to dialogue_comment_update");
        var postData = {
            uid: uid,
            name: $modal.name,
            type_id: $modal.typeId,
            address: $modal.address,
            location_id: $modal.location.LocationId,
            generate_sample_boilers: false,
            generate_sample_data: false
        };
        if ($rootScope.currentUser.Role.RoleId <= 2) {
            postData.show_brand = $modal.showBrand;
            postData.brand_name = $modal.brandName;

            postData.is_super = $modal.isSuper;
            postData.supervisor = $modal.supervisor;
        }

        if (!currentData && (postData.type_id === 1 || postData.type_id === 2)) {
            swal({
                title: "是否为该企业创建示例锅炉？",
                text: "将创建燃煤锅炉、燃气锅炉、生物质锅炉、热水锅炉各一台。注意：所创建的锅炉将录入正式锅炉信息中，如需移除，请在锅炉信息列表中进行删除，或联系平台管理员协助操作。",
                type: "question",
                showCancelButton: true,
                confirmButtonText: "确定",
                cancelButtonText: "取消"
            }).then(function () {
                postData.generate_sample_boilers = true;
                swal({
                    title: "是否为生成示例数据？",
                    text: "示例数据旨在展示平台的数据采集分析流程，不可作为锅炉运行状态检查和故障诊断使用，如需关闭，请联系平台管理员进行关闭。",
                    type: "question",
                    showCancelButton: true,
                    confirmButtonText: "确定",
                    cancelButtonText: "取消"
                }).then(function () {
                    postData.generate_sample_data = true;
                    post(postData);
                }, function () {
                    postData.generate_sample_data = false;
                    post(postData);
                });
            }, function () {
                postData.generate_sample_boilers = false;
                postData.generate_sample_data = false;
                post(postData);
            });
        } else {
            post(postData);
        }

    };

    var post = function (data) {
        $http.post("/organization_update/", data)
            .then(function (res) {
                $rootScope.getOrganizationList();
                organization.refreshDataTables();
                if (data.generate_sample_boilers) {
                    $rootScope.getBoilerList();
                }
                swal({
                    title: "企业信息提交成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close('success');
                    currentData = null;
                });
            }, function (err) {
                swal({
                    title: "企业信息提交失败",
                    text: err.data,
                    type: "error"
                });
            }, function () {
                $rootScope.getOrganizationList();
            });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };
});

// Please note that the close and dismiss bindings are from $uibModalInstance.

angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: '/directives/modal/organization_detail.html',
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