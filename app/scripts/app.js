'use strict';

/**
 * @ngdoc overview
 * @name gohubApp
 * @description
 * # gohubApp
 *
 * Main module of the application.
 */
angular
    .module('gohubApp', [
            'ngResource'
        ], function ($interpolateProvider) {
            $interpolateProvider.startSymbol('[[');
            $interpolateProvider.endSymbol(']]');
        }
    );

