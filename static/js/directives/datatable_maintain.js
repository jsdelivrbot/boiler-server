/**
 * Created by JeremiahYan on 2017/1/8.
 */

boilerAdmin.directive('datatableMaintain', ['$rootScope', function ($rootScope) {
    return {
        restrict: 'E',
        templateUrl: "/directives/datatable_maintain.html",
        // link: function(scope, element, attrs) {
        //     initChartSteamFlow(scope.boilerRuntime);
        // }
    };
}]);