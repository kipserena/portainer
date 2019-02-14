angular.module('extension.storidge').component('storidgeClusterEventsDatatable', {
  templateUrl: 'app/extensions/storidge/components/cluster-events-datatable/storidgeClusterEventsDatatable.html',
  controller: 'GenericDatatableController',
  bindings: {
    titleText: '@',
    titleIcon: '@',
    dataset: '<',
    tableKey: '@',
    orderBy: '@',
    reverseOrder: '<'
  }
});
