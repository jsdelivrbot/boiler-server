/**
 * Created by JeremiahYan on 2017/1/8.
 */

boilerAdmin.directive('datatableHistory', ['$rootScope', function ($rootScope) {
    return {
        restrict: 'E',
        templateUrl: "/directives/datatable_history.html",
        // controller: "BoilerHistoryController",
        // controllerAs: "history"
    };
}]);