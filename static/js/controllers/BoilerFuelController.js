angular.module('BoilerAdmin').controller('BoilerFuelController', function($rootScope, $scope, $http, $timeout, $uibModal, $log, $document, $location, moment, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    boilerFuel = this;

    boilerFuel.isDone = false;

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        boilerFuel.query = $location.search();

        boilerFuel.refreshDataTables();

        // set sidebar closed and body solid layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    boilerFuel.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
        //.withOption('rowCallback', rowCallbackAlarm);

    boilerFuel.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6).notSortable()
    ];

    boilerFuel.refreshDataTables = function () {
        console.info("ready to Get Alarm:", alarm.query);
        $http.get('/fuel_record_list/?boiler=' + alarm.query['boiler'])
            .then(function (res) {
                // $scope.parameters = data;
                var datasource = res.data;

                var num = 0;
                angular.forEach(datasource, function (d, key) {
                    d.num = ++num;

                    d.duration = Math.floor(d.Duration / 1000 / 1000 / 1000);   //sec
                    d.durationDay = Math.floor(d.duration / 60 / 60 / 24);

                    d.duration -= d.durationDay * 24 * 60 * 60;
                    d.durationHour = Math.floor(d.duration / 60 / 60);

                    d.duration -= d.durationHour * 60 * 60;
                    d.durationMin = Math.floor(d.duration / 60);

                    d.duraText = '';
                    d.duraText += d.durationDay > 0 ? d.durationDay + '天' : '';
                    d.duraText += d.durationHour > 0 ? d.durationHour + '小时' : '';
                    d.duraText += d.durationMin + '分';

                    d.statText = boilerFuel.statTexts[d.Status];

                    //d.rtmLen = d.Runtime.length;
                });

                boilerFuel.datasource = datasource;

                boilerFuel.isDone = true;
            });
    };

    function rowCallbackAlarm(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        console.info("rowCallbackAlarm");
        $('td', nRow).unbind('click');
        $('td', nRow).bind('click', function() {
            alarm.editing = false;
            alarm.row = nRow;
            currentData = alarm.datasource[aData[0] - 1];

            alarm.open('lg');
        });
        return nRow;
    }

    boilerFuel.confirm = function (uid) {
        $log.info("alarm.confirm:", uid);
        for (var i = 0; i < boilerFuel.datasource.length; i++) {
            if (boilerFuel.datasource[i].Uid == uid) {
                currentData = boilerFuel.datasource[i];

                $log.info("alarm.confirm GET:", currentData);
                boilerFuel.open('lg');
                break;
            }
        }
    };

    boilerFuel.open = function (size, parentSelector) {
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
            windowClass: 'zindex'
            // resolve: {
            //     items: function () {
            //         return dialogue.items;
            //     }
            // }
        });

        modalInstance.result.then(function (selectedItem) {
            boilerFuel.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };

});

var boilerFuel;
var currentData;

angular.module('BoilerAdmin').controller('ModalAlarmCtrl', function ($uibModalInstance, $http) {
    var $modal = this;
    $modal.editing = false;
    $modal.alarm = currentData;
    $modal.title = "告警详情";

    if (currentData) {
        //$modal.title = "#" + currentData.DialogueId + " " + currentData.Name;
        switch (currentData.Status) {
            case ALARM_STATUS_DEFAULT:
            case ALARM_STATUS_DONE:
                break;
            case ALARM_STATUS_NEW:
            case ALARM_STATUS_PENDING:
            case ALARM_STATUS_REJECTED:
                $modal.feedbackTitle = "确认信息";
                break;
            case ALARM_STATUS_CONFIRMED:
                $modal.feedbackTitle = "审核信息";
                break;
            default:
                break;
        }

        $modal.editing = true;
    }

    /*
     Boiler		*Boiler			`orm:"rel(fk);index"`
     Parameter	*RuntimeParameter	`orm:"rel(fk);index"`
     Runtime		[]*BoilerRuntime	`orm:"reverse(many)"`

     StartDate	time.Time		`orm:"auto_now_add;type(datetime);index"`
     EndDate		time.Time		`orm:"null;type(datetime);index"`

     ConfirmedDate	time.Time		`orm:"null;type(datetime);index"`
     ConfirmedBy	*User			`orm:"rel(fk);null;index"`
     VerifiedDate	time.Time		`orm:"null;type(datetime);index"`
     VerifiedBy	            *User			`orm:"rel(fk);null;index"`

     Feedback	            []*BoilerAlarmFeedback	`orm:"reverse(many)"`

     TriggerRule	        *RuntimeAlarmRule	    `orm:"rel(fk);null;index"`
     AlarmLevel	            int32			        `orm:"index"`
     Status		            int32			        `orm:"index"`
     Priority	            int32			        `orm:"index;default(0)"`
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
            type: 1,
            topic: $modal.topic,
            content: $modal.content
            // attachment: $modal.attachment,
        }).then(function (res) {
            swal({
                title: "提交成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
                alarm.refreshDataTables();
            });
        }, function (err) {
            swal({
                title: "提交失败",
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