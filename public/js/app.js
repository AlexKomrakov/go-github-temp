(function () {
    angular
        .module('gohub', [
            'ngResource'
        ], function($interpolateProvider) {
            $interpolateProvider.startSymbol('[[');
            $interpolateProvider.endSymbol(']]');
        });
})();