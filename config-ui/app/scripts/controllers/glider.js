'use strict';

angular.module('igcstatConfigUiApp')

.controller('GliderCtrl', function($scope, $filter, $http) {
  $scope.gliders = [];

  $scope.saveGlider = function(data, id) {
  };

  // remove user
  $scope.removeGlider = function(index) {
    $scope.gliders.splice(index, 1);
  };

  // add user
  $scope.addGlider = function() {
    $scope.inserted = {
      id: $scope.gliders.length+1,
      igc: '',
      glider: '',
      comment: ''
    };
    $scope.gliders.push($scope.inserted);
  };
});
