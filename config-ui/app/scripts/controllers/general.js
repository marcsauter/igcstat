'use strict';

/**
 * @ngdoc function
 * @name igcstatConfigUiApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the igcstatConfigUiApp
 */
angular.module('igcstatConfigUiApp')

.controller('GeneralCtrl', function ($scope) {
  $scope.user = {
    name: 'awesome user'
  };
});
