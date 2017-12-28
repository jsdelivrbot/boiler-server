angular.module('BoilerAdmin').controller('DialogueController', function($rootScope, $scope, $http, $timeout, $uibModal, $log, $document, moment, settings, DTOptionsBuilder, DTColumnDefBuilder, DTDefaultOptions) {
    dialogue = this;

    dialogue.isDone = false;

    dialogue.exampleDate = moment().hour(8).minute(0).second(0).toDate();

    $scope.$on('$viewContentLoaded', function() {
        // initialize core components
        App.initAjax();

        dialogue.refreshDataTables();

        // set default layout mode
        $rootScope.settings.layout.pageContentWhite = true;
        $rootScope.settings.layout.pageBodySolid = true;
        $rootScope.settings.layout.pageSidebarClosed = false;
    });

    dialogue.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');
        //.withOption('rowCallback', rowCallback);

    dialogue.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2).notSortable(),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5).notSortable()
    ];

    dialogue.refreshDataTables = function () {
        $http.get('/dialogue_list/')
            .then(function (res) {

                var datasource = res.data;

                var num = 0;
                angular.forEach(datasource, function (d, key) {
                    d.num = ++num;
                    var statusTexts = ["默认", "新咨询", "未回复", "已回复", "已关闭"];
                    d.statusText = statusTexts[d.Status];
                    d.name = d.Name;
                    if (d.name.length > 22) {
                        d.name = d.name.substring(0, 20) + "...";
                    }
                    d.title = d.Name;
                    if (d.CreatedBy) {
                        if (d.CreatedBy.Role.RoleId >= 10 && d.CreatedBy.Role.RoleId < 20 && d.CreatedBy.Organization) {
                            d.orgName = d.CreatedBy.Organization.Name;
                        } else {
                            d.orgName = d.CreatedBy.Role.Name;
                        }
                    }
                    //d.orgName = (d.CreatedBy && d.CreatedBy.Organization) ? d.CreatedBy.Organization.Name : " - ";
                    d.username = d.CreatedBy ? d.CreatedBy.Name : " - ";

                    // 2017-02-27T09:21:50+08:00
                    //d.PostDate = moment(d.UpdatedDate, "YYYY-MM-DDTHH:mm:ssZ");
                });

                dialogue.datasource = datasource;

                dialogue.isDone = true;
            });
    };

    function rowCallback(nRow, aData, iDisplayIndex, iDisplayIndexFull) {
        // Unbind first in order to avoid any duplicate handler (see https://github.com/l-lin/angular-datatables/issues/87)
        $('td', nRow).unbind('click');
        $('td', nRow).bind('click', function() {
            dialogue.editing = false;
            dialogue.row = nRow;
            currentData = dialogue.datasource[aData[0] - 1];

            dialogue.open('lg');
        });
        return nRow;
    }

    dialogue.animationsEnabled = true;

    dialogue.new = function () {
        currentData = null;
        dialogue.open('lg');
    };

    dialogue.reply = function (uid) {
        $log.info("dialogue.reply:", uid);
        for (var i = 0; i < dialogue.datasource.length; i++) {
            if (dialogue.datasource[i].Uid == uid) {
                currentData = dialogue.datasource[i];

                $log.info("dialogue.reply GET:", currentData);
                dialogue.open('lg');
                break;
            }
        }
    };

    dialogue.delete = function (uid) {
        swal({
            title: "确认删除该咨询会话？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            $http.post("/dialogue_delete/", {
                uid: uid
            }).then(function (res) {
                swal({
                    title: "咨询会话删除成功",
                    type: "success"
                }).then(function () {
                    dialogue.refreshDataTables();
                });
            }, function (err) {
                swal({
                    title: "删除会话失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

    dialogue.open = function (size, parentSelector) {
        var parentElem = parentSelector ?
            angular.element($document[0].querySelector('.modal-demo ' + parentSelector)) : undefined;
        var modalInstance = $uibModal.open({
            animation: dialogue.animationsEnabled,
            ariaLabelledBy: 'modal-title',
            ariaDescribedBy: 'modal-body',
            templateUrl: '/directives/modal/dialogue_comment.html',
            controller: 'ModalDialogueCtrl',
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
            dialogue.selected = selectedItem;
        }, function () {
            $log.info('Modal dismissed at: ' + new Date());
        });
    };
});

var dialogue;
var currentData;

angular.module('BoilerAdmin').controller('ModalDialogueCtrl', function ($uibModalInstance, $http) {
    var $modal = this;
    $modal.editing = false;
    $modal.dialogue = currentData;

    $modal.title = "开始咨询";
    $modal.commentTitle = "咨询内容";

    if (currentData) {
        $modal.title = "#" + currentData.DialogueId + " " + currentData.Name;
        $modal.commentTitle = "回复内容";
        $modal.editing = true;
    }

    /*
     MyUidObject

     Dialogue	*Dialogue	`orm:"rel(fk);index"`
     From		*User		`orm:"rel(fk);index"`
     To		*User 		`orm:"rel(fk);index;null"`

     Content		string
     Attachment	string
     */

    $modal.ok = function () {
        Ladda.create(document.getElementById('boiler_ok')).start();
        var uid = null;
        var dialogueId = null;
        if (currentData) {
            dialogueId = currentData.Uid;
        }
        //alert("Ready to post to dialogue_comment_update");
        $http.post("/dialogue_comment_update/", {
            uid: uid,
            dialogueId: dialogueId,
            topic: $modal.topic,
            content: $modal.content,
            // attachment: $modal.attachment,
        }).then(function (res) {
            swal({
                title: "提交成功",
                type: "success"
            }).then(function () {
                $uibModalInstance.close('success');
                currentData = null;
                dialogue.refreshDataTables();
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
    templateUrl: '/directives/modal/dialogue_comment.html',
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