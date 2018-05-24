angular.module('BoilerAdmin').controller('AlarmController', function($rootScope, $scope, $http, $timeout,  $filter, $uibModal, $log, $document, $location, moment, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    bAlarm = this;
    $log.warn("Alarm init!");

    bAlarm.isDone = false;

    bAlarm.startFormat = 'YYYY-MM-DD HH:mm';
    bAlarm.endFormat = 'YYYY-MM-DD HH:mm';

    bAlarm.mode = "current";

    bAlarm.dtColumnDefs = [];

    bAlarm.setMode = function (mode) {
        if (bAlarm.mode === mode) {
            return;
        }

        switch (mode) {
            case "current":

                break;
            case "history":
                bAlarm.dtColumnDefs = [
                    DTColumnDefBuilder.newColumnDef(0),
                    DTColumnDefBuilder.newColumnDef(1),
                    DTColumnDefBuilder.newColumnDef(2),
                    DTColumnDefBuilder.newColumnDef(3),
                    DTColumnDefBuilder.newColumnDef(4),
                    DTColumnDefBuilder.newColumnDef(5).notSortable()
                ];
                break;
        }

        bAlarm.mode = mode;
    };

    bAlarm.statusTexts = {
        0: "默认",
        1: "新告警",
        2: "未查阅",
        3: "已查阅",
        4: "驳回",
        5: "已审核",
        10: "已关闭"
    };

    bAlarm.priorityIcons = {
        0: [0],
        1: [0, 1],
        2: [0, 1, 2]
    };

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        bAlarm.query = $location.search();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    bAlarm.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
        //.withOption('rowCallback', rowCallbackAlarm);

    bAlarm.dtColumnDefsCurrent = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6).notSortable()
    ];

    bAlarm.dtColumnDefsHistory = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4)
        // DTColumnDefBuilder.newColumnDef(5).notSortable()
    ];

    $rootScope.$watch('boilerAlarms', function () {
        // $log.warn("$rootScope.$watch.boilerAlarms");
        bAlarm.refreshDataTables();
    });

    bAlarm.declareRefresh = function () {
        $log.warn("Refresh Inline!");
    };

    bAlarm.refreshDataTables = function () {
        //console.info("ready to Get Alarm:", bAlarm.query);
        var datasource = [];

        if (!$rootScope.boilerAlarms) {
            bAlarm.datasource = [];
            console.warn("Alarm Resp Empty:", datasource);
            setTimeout(function () {
                App.stopPageLoading();
            }, 500);
            return;
        }

        if (bAlarm.query['boiler'] && bAlarm.query['boiler'].length > 0) {
            datasource = $filter('filter')($rootScope.boilerAlarms, function (alarm) {
                return (alarm['Boiler__Uid'] === bAlarm.query['boiler']);
            });
        } else {
            // console.warn("Get Alarm Resp:", res);
            datasource = $rootScope.boilerAlarms;
        }

        for (var i = 0; i < datasource.length; i++) {
            var d = datasource[i];
            d.num = i;
            d.isValid = true;
        }

        bAlarm.isDone = true;

        bAlarm.datasource = datasource;
        console.warn("Alarm Resp:", datasource);

        setTimeout(function () {
            App.stopPageLoading();
        }, 800);

        bAlarm.refreshHistory();
    };

    bAlarm.refreshHistory = function() {
        var historyData = [];

        $http.get('/boiler_alarm_history_list/?boiler=' + bAlarm.query['boiler'])
            .then(function (res) {
                console.warn("Get Alarm History List:", res);
                if (res.data) {
                    historyData = res.data;
                }
                for (var i = 0; i < historyData.length; i++) {
                    var d = historyData[i];
                    d.num = i;
                }

                bAlarm.historyData = historyData;

                bAlarm.isDone = true;
            });
    };

    bAlarm.confirm = function (uid) {
        $log.info("bAlarm.confirm:", uid);
        for (var i = 0; i < bAlarm.datasource.length; i++) {
            if (bAlarm.datasource[i].Uid === uid) {
                currentData = bAlarm.datasource[i];

                $log.info("bAlarm.confirm GET:", currentData);
                bAlarm.open('lg');
                break;
            }
        }
    };

    bAlarm.view = function (uid) {
        $log.info("bAlarm.view:", uid);
        for (var i = 0; i < bAlarm.historyData.length; i++) {
            if (bAlarm.historyData[i].Uid === uid) {
                currentData = bAlarm.historyData[i];

                $log.info("bAlarm.view GET:", currentData);
                bAlarm.open('lg');
                break;
            }
        }
    };

    bAlarm.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: true,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/boiler_alarm_feedback.html',
            controller: 'ModalAlarmCtrl',
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
            bAlarm.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

});

var bAlarm;
var currentData;

const ALARM_STATE_DEFAULT = 0;
const ALARM_STATE_NEW = 1;
const ALARM_STATE_PENDING = 2;
const ALARM_STATE_CONFIRMED = 3;
const ALARM_STATE_REJECTED = 4;
const ALARM_STATE_VERIFIED = 5;
const ALARM_STATE_DONE = 10;

angular.module('BoilerAdmin').controller('ModalAlarmCtrl', function ($uibModalInstance, $scope, $http) {
    var $modal = this;
    $modal.editing = false;
    $modal.alarm = {};

    $modal.title = "告警详情";

    $http.post("/boiler_alarm_detail/", {
        uid: currentData.Uid
    }).then(function (res) {
        console.warn("Get Alarm Runtime Resp:", res);
        var alarm = res.data;
        var start = new Date(alarm.StartDate);
        var end = new Date(alarm.EndDate);

        var validTime = 4 * 60 * 60 * 1000;
        var now = new Date();
        var aTime = new Date();
        aTime.setTime(now.getTime() - validTime);
        alarm.isValid = aTime < end;

        alarm.startFormat = 'YYYY-MM-DD HH:mm';
        alarm.endFormat = 'YYYY-MM-DD HH:mm';
        if (end.getYear() === start.getYear()) {
            if (end.getDate() === start.getDate()) {
                alarm.endFormat = 'HH:mm';
            } else {
                alarm.endFormat = 'MM-DD HH:mm';
            }
        }

        $modal.alarm = alarm;
        $scope.alarm = alarm;
        initChartAlarm($scope.alarm);
    }, function (e) {
        console.error("Get Alarm Runtime Failed:", e);
    });

    /*$http.post("/boiler_alarm_feedback_list/", {
        uid: currentData.Uid
    }).then(function (res) {
        console.warn("Get Alarm Feedback Resp:", res);
        $modal.feedbacks = res.data;
    }, function (e) {
        console.error("Get Alarm Runtime Failed:", e);
    });*/


    //$modal.title = "#" + currentData.DialogueId + " " + currentData.Name;
    $modal.feedbackTitle = "备注信息";
    switch ($modal.alarm.State) {
        case ALARM_STATE_DEFAULT:
        case ALARM_STATE_DONE:
            break;
        case ALARM_STATE_NEW:
        case ALARM_STATE_PENDING:
        case ALARM_STATE_REJECTED:
            $modal.feedbackTitle = "确认信息";
            break;
        case ALARM_STATE_CONFIRMED:
            $modal.feedbackTitle = "备注信息";
            break;
        default:
            break;
    }

    $modal.editing = true;

    /*
     Boiler		*Boiler			`orm:"rel(fk);index"`
     Parameter	*RuntimeParameter	`orm:"rel(fk);index"`
     Runtime		[]*BoilerRuntime	`orm:"reverse(many)"`

     StartDate	time.Time		`orm:"auto_now_add;type(datetime);index"`
     EndDate		time.Time		`orm:"null;type(datetime);index"`

     ConfirmedDate	time.Time		`orm:"null;type(datetime);index"`
     ConfirmedBy	*User			`orm:"rel(fk);null;index"`
     VerifiedDate	time.Time		`orm:"null;type(datetime);index"`
     VerifiedBy	*User			`orm:"rel(fk);null;index"`

     Feedback	[]*BoilerAlarmFeedback	`orm:"reverse(many)"`

     TriggerRule	*RuntimeAlarmRule	`orm:"rel(fk);null;index"`
     AlarmLevel	int32			`orm:"index"`
     Status		int32			`orm:"index"`
     Priority	int32			`orm:"index;default(0)"`
     */

    $modal.ok = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();
        var uid = null;
        var alarmId = null;
        if (currentData) {
            alarmId = currentData.Uid;
        }
        //alert("Ready to post to dialogue_comment_update");
        $http.post("/boiler_alarm_update/", {
            uid: uid,
            alarm_id: alarmId,
            state: 1,
            topic: $modal.topic,
            content: $modal.content,
            // attachment: $modal.attachment,
        }).then(function (res) {
            currentData.State = ALARM_STATE_CONFIRMED;
            bAlarm.refreshDataTables();
            Ladda.create(document.getElementById('boiler_ok')).stop();
            swal({
                title: "提交成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
            });
        }, function (err) {
            Ladda.create(document.getElementById('boiler_ok')).stop();
            swal({
                title: "提交失败",
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
    templateUrl: '/directives/modal/boiler_alarm_feedback.html',
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