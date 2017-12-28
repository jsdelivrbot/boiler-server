/**
 * Created by JeremiahYan on 2017/1/8.
 */

boilerAdmin.directive('datatableAlarm', ['$rootScope', function ($rootScope) {
    return {
        restrict: 'E',
        templateUrl: "/directives/datatable_alarm.html"
    };
}]);

boilerAdmin.directive('datatableAlarmHistory', ['$rootScope', function ($rootScope) {
    return {
        restrict: 'E',
        templateUrl: "/directives/datatable_alarm_history.html"
    };
}]);