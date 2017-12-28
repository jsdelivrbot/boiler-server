/**
 * Created by JeremiahYan on 2017/2/5.
 */

boilerAdmin.directive('locationSelect', ['$rootScope', function ($rootScope) {
    return {
        restrict: 'E',
        // require: 'ngModel',
        templateUrl: "/directives/location-select.html",
        // link: function () {
        //     $rootScope.province = 0;
        //     $rootScope.city = null;
        //     $rootScope.$watch('province', function (newValue, oldValue) {
        //         alert(newValue + ", " + oldValue);
        //     });
        //
        //     $rootScope.change = function () {
        //         alert("Province: " + $rootScope.province);
        //     }
        // }
        // link: function (scope, element, attrs, ngModel) {
        //     attrs.$observe('ngModel', function (value) {
        //         scope.$watch(value, function (newValue) {
        //             alert(newValue);
        //         })
        //     });
        // }
        // link: function (scope, element, attrs, ngModel) {
        //     ngModel.$render = function () {
        //         var newValue = ngModel.$viewValue;
        //         console.log(newValue)
        //         alert(newValue);
        //     };
        //
        // }
    };
}]);

var local;

boilerAdmin.controller('LocationController', ['$scope', '$rootScope', '$http', function($scope, $rootScope, $http) {
    local = this;

    $scope.$on('$viewContentLoaded', function() {
        local.province = null;
        local.city = null;
        local.region = null;
    });

    local.changeProvince = function () {
        local.location = local.province;
        //alert(local.location.Name);
    };

    local.changeCity = function () {
        local.location = local.city;
        //alert(local.location.Name);
    };

    local.changeRegion = function () {
        local.location = local.region;
        //alert(local.location.Name);
    };
}]);