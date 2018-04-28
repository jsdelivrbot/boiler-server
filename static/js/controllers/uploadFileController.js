angular.module('BoilerAdmin').controller("uploadFileCtrl",function ($rootScope,$scope,$uibModal,$http,$document,settings,DTOptionsBuilder, DTColumnDefBuilder) {
   var upload = this;

    App.initAjax();

    // dialogue.refreshDataTables();

    // set default layout mode
    $rootScope.settings.layout.pageContentWhite = true;
    $rootScope.settings.layout.pageBodySolid = true;
    $rootScope.settings.layout.pageSidebarClosed = false;

    upload.dtOptions = DTOptionsBuilder.newOptions()
        .withPaginationType('full_numbers');

    upload.dtColumnDefs = [
        DTColumnDefBuilder.newColumnDef(0),
        DTColumnDefBuilder.newColumnDef(1),
        DTColumnDefBuilder.newColumnDef(2),
        DTColumnDefBuilder.newColumnDef(3),
        DTColumnDefBuilder.newColumnDef(4),
        DTColumnDefBuilder.newColumnDef(5),
        DTColumnDefBuilder.newColumnDef(6).notSortable()
    ];

    upload.refresh = function () {
      $http.get("/bin_list").then(function (res) {
          upload.datasource = res.data;
          for(var i = 0; i<upload.datasource.length; i++){
              upload.datasource[i].num = i+1;
          }

      },function (err) {

      })
    };
    upload.refresh();

    upload.new = function () {
        var modalInstance = $uibModal.open({
            templateUrl: 'addFile.html',
            controller: 'ModalFileUploadCtrl',
            size: "lg",
            windowClass: 'zindex',
        });


        modalInstance.result.then(function (selectedItem) {
            $scope.selected = selectedItem;
        }, function () {

        });
    }


    upload.delete = function (name) {
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
            $http.post("/bin_delete", {
                name: name
            }).then(function (res) {
                swal({
                    title: "文件删除成功",
                    type: "success"
                }).then(function () {
                    upload.refresh();
                });
            }, function (err) {
                swal({
                    title: "删除文件失败",
                    text: err.data,
                    type: "error"
                });
            });
        });
    };

})


angular.module('BoilerAdmin').controller('ModalFileUploadCtrl', function ($scope,$rootScope, $uibModalInstance) {

    console.log($rootScope.organizations);
    $scope.organizations = $rootScope.organizations;

    $scope.org = function(organization){
        $scope.orgUid = organization.Uid;
        console.log($scope.orgUid);
    }

    $scope.ok = function () {
        $uibModalInstance.close();
    };

    $scope.cancel = function () {
        $uibModalInstance.dismiss('cancel');
    };
});
