angular.module('BoilerAdmin').controller("templateCtrl",function ($rootScope,$scope,$uibModal,$document,settings,DTOptionsBuilder, DTColumnDefBuilder) {
    var template = this;

    App.initAjax();

    // dialogue.refreshDataTables();

    // set default layout mode
    $rootScope.settings.layout.pageContentWhite = true;
    $rootScope.settings.layout.pageBodySolid = true;
    $rootScope.settings.layout.pageSidebarClosed = false;

    template.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');

    template.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3).notSortable()
    ];

    template.datasource=[
        {
            num:1,
            name:"通用模板一",
            enterprise:"锅炉制造厂"
        },
        {
            num:2,
            name:"通用模板二",
            enterprise:"二号锅炉制造厂"
        }
    ];

    template.new = function () {
        var modalInstance = $uibModal.open({
            templateUrl: 'editTemplate.html',
            controller: 'ModalEditTemplateCtrl',
            size: "lg",
            windowClass: 'zindex',
        });


        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    };

    template.edit = function () {
        var modalInstance = $uibModal.open({
            templateUrl: 'editTemplate.html',
            controller: 'ModalEditTemplateCtrl',
            size: "lg",
            windowClass: 'zindex',
        });


        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    }


    template.delete = function (uid) {
        swal({
            title: "确认删除该文件？",
            text: "注意：删除后将无法恢复",
            type: "warning",
            showCancelButton: true,
            //confirmButtonClass: "btn-danger",
            confirmButtonColor: "#d33",
            cancelButtonText: "取消",
            confirmButtonText: "删除",
            closeOnConfirm: false
        }).then(function () {
            // $http.post("/dialogue_delete/", {
            //     uid: uid
            // }).then(function (res) {
            //     swal({
            //         title: "文件删除成功",
            //         type: "success"
            //     }).then(function () {
            //         dialogue.refreshDataTables();
            //     });
            // }, function (err) {
            //     swal({
            //         title: "删除文件失败",
            //         text: err.data,
            //         type: "error"
            //     });
            // });
        });
    };

})


angular.module('BoilerAdmin').controller('ModalEditTemplateCtrl', function ($scope, $uibModalInstance) {

    $scope.ok = function () {
        $uibModalInstance.close();
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});
