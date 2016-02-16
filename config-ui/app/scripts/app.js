'use strict';

/**
 * @ngdoc overview
 * @name igcstatConfigUiApp
 * @description
 * # igcstatConfigUiApp
 *
 * Main module of the application.
 */
angular.module('igcstatConfigUiApp', [
  'ngAnimate',
  'ngCookies',
  'ngResource',
  'ngSanitize',
  'ngTouch',
  'ui.router',
  'xeditable'
])

.config(function($stateProvider, $urlRouterProvider) {
  //
  // For any unmatched url, redirect to /state1
  $urlRouterProvider.otherwise("/general");
  //
  // Now set up the states
  $stateProvider
    .state('general', {
      url: "/general",
      templateUrl: "views/general.html"
    })
    .state('flight', {
      url: "/flight",
      templateUrl: "views/flight.html"
    })
    .state('glider', {
      url: "/glider",
      templateUrl: "views/glider.html"
    });
})

.run(function (editableOptions) {
  editableOptions.theme = 'bs3'; // bootstrap3 theme. Can be also 'bs2', 'default'
});
