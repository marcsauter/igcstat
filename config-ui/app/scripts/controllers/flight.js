'use strict';

angular.module('igcstatConfigUiApp')

.controller('FlightCtrl', function($scope, $filter, $http) {
  $scope.flights = [];

  $scope.saveFlight = function(data, id) {
  };

  // remove user
  $scope.removeFlight = function(index) {
    $scope.flights.splice(index, 1);
  };

  // add user
  $scope.addFlight = function() {
    $scope.inserted = {
      id: $scope.flights.length+1,
      date: '',
      takeOffTime: '',
      takeOffSite: '', 
      landingTime: '',
      landingSite: '',
      glider: '',
      comment: ''
    };
    $scope.flights.push($scope.inserted);
  };
});
