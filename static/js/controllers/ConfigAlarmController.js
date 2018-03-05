angular.module('BoilerAdmin').controller('ConfigAlarmController', function($rootScope, $scope, $http, $timeout, $uibModal, $log, $document, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    confAlarm = this;
    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        confAlarm.refreshDataTables();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    confAlarm.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers')
        .withOption('rowCallback', rowCallback);

    confAlarm.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6),
        DTColumnDefBuilder.newColumnDef(7),
        DTColumnDefBuilder.newColumnDef(8),
        DTColumnDefBuilder.newColumnDef(9).notSortable()
    ];

    confAlarm.refreshDataTables = function () {
        App.startPageLoading({message: '正在加载数据...'});
        $http.get('/alarm_rule_list/')
            .then(function (res) {
                // $scope.parameters = data;
                var datasource = res.data;

                var num = 0;
                angular.forEach(datasource, function (d, key) {
                    d.num = ++num;
                    var priorityTexts = ["低", "中", "高"];
                    d.priortyText = priorityTexts[d.Priority];
                    d.formName = d.BoilerForm ? d.BoilerForm.Name : " - ";
                    d.mediumName = d.BoilerMedium ? d.BoilerMedium.Name.substring(0, d.BoilerMedium.Name.length - 2) : " - ";
                    d.fuelName = d.BoilerFuelType ? d.BoilerFuelType.Name : " - ";
                    d.warning = d.Warning > d.Normal ? " ＞ " + d.Warning : " ＜ " + d.Warning;
                    //d.danger = d.Danger > 0 ? d.Danger : " - ";
                    d.capacity = " 不限 ";
                    if (d.BoilerCapacityMax > d.BoilerCapacityMin) {
                        d.capacity = d.BoilerCapacityMin + " - " + d.BoilerCapacityMax;
                    } else if (d.BoilerCapacityMin > 0) {
                        d.capacity = d.BoilerCapacityMin;
                    }
                });

                confAlarm.datasource = datasource;
                setTimeout(function () {
                    App.stopPageLoading();
                }, 500);
            });
    };

    $scope.$on('modal.closing', function(event, reason, closed) {
        console.log('modal.closing: ' + (closed ? 'close' : 'dismiss') + '(' + reason + ')');
        var message = "You are about to leave the edit view. Uncaught reason. Are you sure?";
        switch (reason){
            // clicked outside
            case "backdrop click":
                message = "Any changes will be lost, are you sure?";
                break;

            // cancel button
            case "cancel":
                message = "Any changes will be lost, are you sure?";
                break;

            // escape key
            case "escape key press":
                message = "Any changes will be lost, are you sure?";
                break;
        }
        if (!confirm(message)) {
            event.preventDefault();
        }
    });

    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        $('td', nRow).unbind('click');
        $('td', nRow).bind('click', function() {
            confAlarm.editing = false;
            confAlarm.row = nRow;
            currentData = confAlarm.datasource[aData[0] - 1];

            confAlarm.open();

            // $scope.$apply(function() {
            //     someClickHandler(confAlarm.currentData);
            // });
        });
        return nRow;
    }
    // confAlarm.items = ['item1', 'item2', 'item3'];

    confAlarm.animationsEnabled = true;

    confAlarm.new = function () {
        currentData = null;
        confAlarm.open();
    };

    confAlarm.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: confAlarm.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/alarm_rule_config.html',
            controller: 'ModalAlarmRuleCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            // resolve: {
            //     items: function () {
            //         return confAlarm.items;
            //     }
            // }
        });

        modalInstance.result.then(function (selectedItem) {
            confAlarm.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

    confAlarm.toggleAnimation = function () {
        confAlarm.animationsEnabled = !confAlarm.animationsEnabled;
    };
});

var confAlarm;
var currentData;

angular.module('BoilerAdmin').controller('ModalAlarmRuleCtrl', function ($uibModalInstance, $http) {
    var $modal = this;
    $modal.editing = false;
    // $modal.items = items;
    $modal.title = "新建告警规则";

    $modal.boilerFormId = 0;
    $modal.boilerMediumId = 0;
    $modal.boilerFuelTypeId = 0;

    $modal.delay = 10;

    $modal.priority = 1;

    if (currentData) {
        $modal.editing = true;
        $modal.title = "编辑告警规则";

        $modal.paramId = currentData.Parameter.Id;
        $modal.boilerFormId = currentData.BoilerForm.Id;
        $modal.boilerMediumId = currentData.BoilerMedium.Id;
        $modal.boilerFuelTypeId = currentData.BoilerFuelType.Id;
        $modal.boilerCapacityMin = currentData.BoilerCapacityMin;
        $modal.boilerCapacityMax = currentData.BoilerCapacityMax;

        $modal.normalValue = currentData.Normal;
        $modal.warningValue = currentData.Warning;

        $modal.delay = currentData.Delay;
        $modal.priority = currentData.Priority;

        $modal.description = currentData.Description;
    }

    /*
     BoilerForm		*BoilerTypeForm			`orm:"rel(fk);null;index"`
     BoilerMedium		*BoilerMedium			`orm:"rel(fk);null;index"`
     BoilerFuelType		*FuelType			`orm:"rel(fk);null;index"`
     BoilerCapacityMin	int32				`orm:"index"`
     BoilerCapacityMax	int32				`orm:"index"`
     Normal			float32				//基准值
     Warning			float32				//告警值
     Danger			float32				//危险值
     Delay			int64				//延迟报警时间，单位分
     Priority		int32				`orm:"index;default(0)"`
     Scope			int32
     */

    $modal.delete = function () {
        var uid = null;
        if (currentData) {
            uid = currentData.Uid;
        }
        swal({
            title: "确认删除该参数？\n" + currentData.Name,
            text: "注意：删除后将无法恢复。",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/alarm_rule_delete/", {
                uid: uid
            }).then(function (res) {
                swal({
                    title: "参数删除成功",
                    type: "success"
                }).then(function () {
                    $uibModalInstance.close();
                    currentData = null;
                    confAlarm.refreshDataTables();
                });
            }, function (err) {
                swal({
                    title: "参数删除失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

    $modal.ok = function () {
        // Ladda.create(document.getElementById('boiler_ok')).start();
        var uid = null;
        if (currentData) {
            uid = currentData.Uid;
        }
        $http.post("/alarm_rule_update/", {
            uid: uid,
            paramId: $modal.paramId,
            boilerFormId: $modal.boilerFormId,
            boilerMediumId: $modal.boilerMediumId,
            boilerFuelTypeId: $modal.boilerFuelTypeId,
            boilerCapacityMin: $modal.boilerCapacityMin,
            boilerCapacityMax: $modal.boilerCapacityMax,

            normalValue: parseFloat($modal.normalValue),
            warningValue: parseFloat($modal.warningValue),
            delay: parseInt($modal.delay),
            priority: $modal.priority,

            description: $modal.description
        }).then(function (res) {
            // Ladda.create(document.getElementById('boiler_ok')).stop();
            swal({
                title: "告警规则更新成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
                confAlarm.refreshDataTables();
            });
        }, function (err) {
            // Ladda.create(document.getElementById('boiler_ok')).stop();
            swal({
                title: "告警规则更新失败",
                text: err.data,
                type: "error"
            });
        });
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };
});

// Please note that the close and dismiss bindings are from $uibModalInstance.

angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: '/directives/modal/alarm_rule_config.html',
    bindings: {
        resolve: '<',
        close: '&',
        dismiss: '&'
    },
    controller: function () {
        var $ctrl = this;

        // $ctrl.$onInit = function () {
        //     $ctrl.items = $ctrl.resolve.items;
        //     $ctrl.selected = {
        //         item: $ctrl.items[0]
        //     };
        // };

        $ctrl.ok = function () {
            $ctrl.close({$value: 'success'});
        };

        $ctrl.cancel = function () {
            $ctrl.dismiss({$value: 'cancel'});
        };
    }
});