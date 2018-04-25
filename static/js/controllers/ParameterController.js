angular.module('BoilerAdmin').controller('ParameterController', function ($rootScope, $scope, $http, $timeout, $uibModal, $log, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    var bParameter = this;

    bParameter.category = PARAMETER_CATEGORY_ANALOG;
    bParameter.categoryList = [
        {
            Id: PARAMETER_CATEGORY_UNDEFINED,
            Name: '请选择...'
        },
        {
            Id: PARAMETER_CATEGORY_ANALOG,
            Name: '模拟量'
        },
        {
            Id: PARAMETER_CATEGORY_SWITCH,
            Name: '开关量'
        },
        // {
        //     Id: PARAMETER_CATEGORY_CALCULATE,
        //     Name: '计算量'
        // },
        {
            Id: PARAMETER_CATEGORY_STATUS,
            Name: '状态量'
        }
    ];

    var currentData;
    var isNew = false;
    var editing = false;

    $rootScope.$watch('parameters', function () {
        // $log.warn("$rootScope.$watch.boilers: ", $rootScope.boilers);
        if (!$rootScope.parameters || typeof $rootScope.parameters !== 'object') {
            return;
        }

        bParameter.refreshDataTables(bParameter.category);
    });

    bParameter.refreshDataTables = function (category) {
        bParameter.category = category;

        var datasource = [];
        $http.get("/runtime_parameter_issued_list").then(function (res) {
            var parameters = res.data;
            for (var i = 0; i < parameters.length; i++) {
                var param = parameters[i];

                if (param.Parameter.Category.Id === category) {
                    datasource.push(param);
                }

                bParameter.datasource = datasource;
            }

            setTimeout(function () {
                App.stopPageLoading();
            }, 1500);
        },function (err) {
            setTimeout(function () {
                App.stopPageLoading();
            }, 1500);
        });


    };

    bParameter.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
    //.withOption('rowCallback', rowCallback);

    bParameter.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6).notSortable(),
        DTColumnDefBuilder.newColumnDef(7).notSortable()
    ];

    bParameter.new = function () {
        currentData = {
            Parameter:{
                Category: {
                    Id: PARAMETER_CATEGORY_UNDEFINED
                }
            }
        };
        isNew = true;
        editing = true;

        bParameter.open('lg');
    };

    bParameter.edit = function (data) {
        currentData = data;
        isNew = false;
        editing = true;

        bParameter.open('lg');
    };

    bParameter.view = function (data) {
        currentData = data;
        isNew = false;
        editing = false;

        bParameter.open('lg');
    };

    bParameter.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/parameter_config.html',
            controller: 'ModalParameterCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                parameter: function () {
                    return bParameter;
                },
                isNew: function () {
                    return isNew;
                },
                editing: function () {
                    return editing;
                },
                currentData: function () {
                    return currentData;
                },
                categoryList: function () {
                    return bParameter.categoryList;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            // bParameter.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    $scope.$on('$viewContentLoaded', function () {
        // initialize core components
        App.initAjax();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });
});

angular.module('BoilerAdmin').controller('ModalParameterCtrl', function ($uibModalInstance, $uibModal, $rootScope, $http, $log, parameter, isNew, editing, currentData, categoryList) {
    var $modal = this;
    $modal.data = currentData;
    $modal.categoryList = categoryList;

    $modal.isNew = isNew;
    $modal.editing = editing;

    $modal.title = isNew ? "创建参数" : editing ? "配置参数" : "查看参数";

    /*
    $modal.initCurrent = function () {
        if (currentData) {

            if (!currentData.Boilers) {
                currentData.Boilers = [];
            }

            for (var i = 0; i < 8; i++) {
                if (i < currentData.Boilers.length) {
                    var boiler = currentData.Boilers[i];
                    boiler.num = boiler.TerminalSetId;
                    boiler.hasDev = true;
                    $modal.sets.push(boiler);
                } else {
                    $modal.sets.push({
                        num: i + 1,
                        Name: "未配置",
                        hasDev: false
                    });
                }
            }

            $modal.deviceCount = currentData.Boilers.length;
        }
    };

    $modal.initCurrent();
    */
    $modal.categoryChanged = function () {
        var cateId = $modal.data.Parameter.Category.Id;
        if (cateId <= 0) {
            $modal.data.Parameter.ParamId = 0;
            $modal.data.Parameter.Id = 0;
        }

        if (cateId === PARAMETER_CATEGORY_SWITCH || cateId === PARAMETER_CATEGORY_STATUS) {
            $modal.data.Parameter.Scale = 1;
            $modal.data.Parameter.Fix = 0;
            $modal.data.Parameter.Unit = "";
            $modal.data.Parameter.Length = 1;
        } else {
            $modal.data.Parameter.Fix = 2;
            $modal.data.Parameter.Length = 2;
        }

        var paramId = 100;
        for (var i = 0; i < $rootScope.parameters.length; i++) {
            var p = $rootScope.parameters[i];
            if (p.Category.Id === cateId && p.ParamId >= paramId) {
                paramId = p.ParamId + 1;
            }
        }

        $modal.data.Parameter.ParamId = paramId;
        $modal.data.Parameter.Id = parseInt(cateId + '' + paramId);
    };

    $modal.commit = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();

        $http.post("/runtime_parameter_update/", $modal.data)
            .then(function (res) {
                swal({
                    title: "参数更新成功",
                    text: "您可以在 终端管理 -> 通道配置 中进行该参数的通道配置",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close('success');
                    $rootScope.getParameterList();
                });
            }, function (err) {
                swal({
                    title: "参数更新失败",
                    text: err.data,
                    type: "error"
                });
            });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.delete = function () {
        swal({
            title: "确认删除该参数？",
            text: "注意：删除后将无法恢复，且和此参数相关的数据通道将废止使用。",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/runtime_parameter_delete/", {
                Id: $modal.data.Id
            }).then(function (res) {
                swal({
                    title: "参数删除成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close('success');
                    $rootScope.getParameterList();
                });
            }, function (err) {
                swal({
                    title: "删除参数失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };
});

const PARAMETER_CATEGORY_UNDEFINED  = 0;
const PARAMETER_CATEGORY_ANALOG     = 10;
const PARAMETER_CATEGORY_SWITCH     = 11;
const PARAMETER_CATEGORY_CALCULATE  = 12;
const PARAMETER_CATEGORY_STATUS     = 13;