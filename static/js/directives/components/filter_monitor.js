/**
 * Created by JeremiahYan on 2017/5/12.
 */

boilerAdmin.directive('filterMonitor', ['$rootScope', function ($rootScope) {
    return {
        restrict: 'E',
        templateUrl: "/directives/components/filter_monitor.html",
        // link: initFlotChartSmoke($rootScope.boilerRuntime)
    };
}]);