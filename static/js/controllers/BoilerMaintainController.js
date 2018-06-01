angular.module('BoilerAdmin').controller('BoilerMaintainController', function($rootScope, $scope, $http, $timeout, $uibModal, $log, $document, $location, settings, moment, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    bMaintain = this;

    bMaintain.isDone = false;

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        bMaintain.query = $location.search();
        bMaintain.refreshDataTables();

        // set default layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    bMaintain.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
        //.withOption('rowCallback', rowCallbackMaintain);

    bMaintain.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4).notSortable(),
        DTColumnDefBuilder.newColumnDef(5).notSortable()
    ];

    bMaintain.getBoilers = function () {
        bMaintain.boilers = [];
        for (var i = 0; i < $rootScope.boilers.length; i++) {
            var d = $rootScope.boilers[i];
            d.name = d.Name;
            if (d.Enterprise) {
                d.ent = d.Enterprise.Name;
            } else {
                d.ent = "";
            }

            bMaintain.boilers.push(d);
        }
    };

    bMaintain.getBoilers();

    $rootScope.$watch('boilers', function () {
        bMaintain.getBoilers();
    });

    var stringToInts = function (str) {
        var ar = str.split(',');
        for(var i = 0; i < ar.length; i++) {
            ar[i] = +ar[i];
        }

        return ar;
    };

    bMaintain.refreshDataTables = function () {
        var p = $location.search();
        console.info("ready to Get Maintain:", p);
        $http.get('/boiler_maintain_list/?boiler=' + p['boiler'])
            .then(function (res) {
                console.warn("Get Maintain Resp:", res);
                var datasource = res.data;

                var num = 0;
                angular.forEach(datasource, function (d, key) {
                    d.num = ++num;
                    d.summary = d.Content;
                    if (d.summary.length > 26) {
                        d.summary = d.summary.substring(0, 26) + "...";
                    }
                    d.status = {
                        burner: stringToInts(d.Burner),                 // 燃烧器
                        importGrate: stringToInts(d.ImportGrate),       // 进料及炉排
                        waterSoftener: stringToInts(d.WaterSoftener), 	// 软水器
                        waterPump: stringToInts(d.WaterPump), 	        // 水泵
                        boilerBody: stringToInts(d.BoilerBody),         // 锅炉本体
                        energySaver: stringToInts(d.EnergySaver),	    // 节能器
                        airPreHeater: stringToInts(d.AirPreHeater),	    // 空预器
                        dustCatcher: stringToInts(d.DustCatcher),	    // 除尘器
                        draughtFan: stringToInts(d.DraughtFan)	        // 引风机
                    };
                });

                bMaintain.datasource = datasource;

                bMaintain.isDone = true;

                if (p["boiler"]) {
                    bMaintain.currentBoilerId = p["boiler"];
                }

            }, function (err) {

            });
    };

    function rowCallbackMaintain(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        $('td', nRow).unbind('click');
        $('td', nRow).bind('click', function() {
            bMaintain.editing = false;
            bMaintain.row = nRow;
            currentData = bMaintain.datasource[aData[0] - 1];

            bMaintain.open('lg');

            $scope.$apply(function() {
                //someClickHandler(bMaintain.currentData);
            });
        });
        return nRow;
    }

    bMaintain.animationsEnabled = true;

    bMaintain.new = function () {
        currentData = null;
        bMaintain.open('lg');
    };

    bMaintain.edit = function (uid) {
        $log.info("bMaintain.edit:", uid);
        for (var i = 0; i < bMaintain.datasource.length; i++) {
            if (bMaintain.datasource[i].Uid == uid) {
                currentData = bMaintain.datasource[i];
                //$log.info("bMaintain.edit GET:", currentData);
                bMaintain.open('lg');
                break;
            }
        }
    };

    bMaintain.delete = function (uid) {
        swal({
            title: "确认删除该记录？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/boiler_maintain_delete/", {
                uid: uid
            }).then(function (res) {
                swal({
                    title: "维保记录删除成功",
                    type: "success"
                }).then(function () {
                    bMaintain.refreshDataTables();
                });
            }, function (err) {
                swal({
                    title: "删除记录失败",
                    text: err.data,
                    type: "error"
                });
            });
        });

    };

    bMaintain.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: bMaintain.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/boiler_maintain_detail.html',
            controller: 'ModalMaintainCtrl',
            controllerAs: '$modal',
            size: size,
            appendTo: parentElem,
            windowClass: 'zindex',
            resolve: {
                currentBoilerId: function () {
                    return bMaintain.currentBoilerId;
                }
            }
        });

        modalInstance.result.then(function (selectedItem) {
            bMaintain.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };
});

var bMaintain;
var currentData;

angular.module('BoilerAdmin').controller('ModalMaintainCtrl', function ($uibModalInstance, $scope, $http, currentBoilerId) {
    var $modal = this;
    $modal.currentData = currentData;
    $modal.hasCurrentBoiler = false;

    for (var i in bMaintain.boilers) {
        var boiler = bMaintain.boilers[i];
        if (boiler.Uid === currentBoilerId) {
            $modal.boilers = [boiler];
            $modal.boilerId = currentBoilerId;
            $modal.hasCurrentBoiler = true;
            break;
        }
    }

    if (!$modal.hasCurrentBoiler) {
        $modal.boilers = bMaintain.boilers;
    }

    $modal.today = function() {
        $modal.inspectDate = new Date();
    };

    $modal.maintainDetail = {
        burner: [0, 0, 0, 0, 0, 0, 0],  // 燃烧器
        importGrate: [0, 0, 0, 0, 0, 0, 0],	// 进料及炉排
        waterSoftener: [0, 0, 0],	// 软水器
        waterPump: [0, 0, 0, 0],	// 水泵
        boilerBody: [0, 0, 0, 0, 0, 0],// 锅炉本体
        energySaver: [0, 0, 0],	// 节能器
        airPreHeater: [0, 0, 0],	// 空预器
        dustCatcher: [0, 0, 0],	// 除尘器
        draughtFan: [0, 0, 0],	// 引风机
    };

    if (currentData) {
        $modal.title = "编辑参数";
        $modal.editing = true;

        $modal.boilerId = currentData.Boiler.Uid;
        $modal.inspectDate = new Date(currentData.InspectDate);
        $modal.content = currentData.Content;
        $modal.maintainDetail = currentData.status;
    } else {
        $modal.today();
    }

    /*
     MyUidObject

     Boiler		*Boiler			`orm:"rel(fk);null;index"`

     InspectDate	time.Time		`orm:"type(datetime);index"`
     Inspector	string
     Content		string
     Attachment	string
     IsDone		bool


     */
    $modal.detailindex = 1;
    $modal.setIndex = function (idx) {
        $modal.detailindex = idx;
    };

    $modal.ok = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();
        var uid = null;
        $modal.burner_jsop = $modal.maintainDetail.burner.join(",");
        $modal.import_grate_jsop = $modal.maintainDetail.importGrate.join(",");
        $modal.water_softener_jsop = $modal.maintainDetail.waterSoftener.join(",");
        $modal.water_pump_jsop = $modal.maintainDetail.waterPump.join(",");
        $modal.boiler_body_jsop = $modal.maintainDetail.boilerBody.join(",");
        $modal.energy_saver_jsop = $modal.maintainDetail.energySaver.join(",");
        $modal.air_pre_heater_jsop = $modal.maintainDetail.airPreHeater.join(",");
        $modal.dust_catcher_jsop = $modal.maintainDetail.dustCatcher.join(",");
        $modal.draught_fan_jsop = $modal.maintainDetail.draughtFan.join(",");
        //alert("Ready to post to maintain_comment_update");
        $http.post("/boiler_maintain_update/", {
            uid: uid,
            boiler_id: $modal.boilerId,
            inspect_date: $modal.inspectDate,
            topic: $modal.topic,
            content: $modal.content,
            attachment: $modal.attachment,
            is_done: $modal.isDone,

            burner: $modal.burner_jsop,
            import_grate: $modal.import_grate_jsop,
            water_softener: $modal.water_softener_jsop,
            water_pump: $modal.water_pump_jsop,
            boiler_body: $modal.boiler_body_jsop,
            energy_saver: $modal.energy_saver_jsop,
            air_pre_heater: $modal.air_pre_heater_jsop,
            dust_catcher: $modal.dust_catcher_jsop,
            draught_fan: $modal.draught_fan_jsop,

        }).then(function (res) {
            swal({
                title: "锅炉维保记录提交成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
                bMaintain.refreshDataTables();
            });
        }, function (err) {
            swal({
                title: "锅炉维保记录提交失败",
                text: err.data,
                type: "error"
            });
        });
        Ladda.create(document.getElementById('boiler_ok')).stop();
    };

    $modal.cancel = function () {
        $uibModalInstance.dismiss('cancel');

        currentData = null;
    };

    $modal.today = function() {
        $modal.inspectDate = new Date();
    };

    $scope.clear = function() {
        $modal.inspectDate = null;
    };

    $scope.inlineOptions = {
        customClass: getDayClass,
        minDate: new Date(),
        showWeeks: true
    };

    $scope.dateOptions = {
        dateDisabled: disabled,
        formatYear: 'yy',
        maxDate: new Date(2020, 5, 22),
        minDate: new Date(),
        startingDay: 1
    };

    // Disable weekend selection
    function disabled(data) {
        var date = data.date,
            mode = data.mode;
        return mode === 'day' && (date.getDay() === 0 || date.getDay() === 6);
    }

    $scope.toggleMin = function() {
        $scope.inlineOptions.minDate = $scope.inlineOptions.minDate ? null : new Date();
        $scope.dateOptions.minDate = $scope.inlineOptions.minDate;
    };

    $scope.toggleMin();

    $scope.open2 = function() {
        $scope.popup2.opened = true;
    };

    $scope.setDate = function(year, month, day) {
        $scope.dt = new Date(year, month, day);
    };

    $scope.formats = ['dd-MMMM-yyyy', 'yyyy/MM/dd', 'dd.MM.yyyy', 'shortDate'];
    $scope.format = $scope.formats[0];
    $scope.altInputFormats = ['M!/d!/yyyy'];

    $scope.popup2 = {
        opened: false
    };

    var tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    var afterTomorrow = new Date();
    afterTomorrow.setDate(tomorrow.getDate() + 1);
    $scope.events = [
        {
            date: tomorrow,
            status: 'full'
        },
        {
            date: afterTomorrow,
            status: 'partially'
        }
    ];

    function getDayClass(data) {
        var date = data.date,
            mode = data.mode;
        if (mode === 'day') {
            var dayToCheck = new Date(date).setHours(0,0,0,0);

            for (var i = 0; i < $scope.events.length; i++) {
                var currentDay = new Date($scope.events[i].date).setHours(0,0,0,0);

                if (dayToCheck === currentDay) {
                    return $scope.events[i].status;
                }
            }
        }

        return '';
    }
});

// Please note that the close and dismiss bindings are from $uibModalInstance.

angular.module('BoilerAdmin').component('modalComponent', {
    templateUrl: '/directives/modal/boiler_maintain_detail.html',
    bindings: {
        resolve: '<',
        close: '&',
        dismiss: '&'
    },
    controller: function () {
        var $ctrl = this;

        $ctrl.$onInit = function () {
            // $ctrl.items = $ctrl.resolve.items;
            $ctrl.selected = {
                // item: $ctrl.items[0]
            };
        };

        $ctrl.ok = function () {
            // $ctrl.close({$value: $ctrl.selected.item});
        };

        $ctrl.cancel = function () {
            $ctrl.dismiss({$value: 'cancel'});
        };
    }
});