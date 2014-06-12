/*
 * Copyright (c) 2014 Jason Ish
 * All rights reserved.
 */

/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

app.directive("eventDetail", function() {

    directive = {
        restrict: "A",
        templateUrl: "views/detail.html",
        scope: {
            hit: "=event"
        }
    };

    directive.controller = function($scope, Keyboard, Util) {

        $scope.Util = Util;

        $scope.$on("$destroy", function () {
            Keyboard.resetScope($scope);
        });

        Keyboard.scopeBind($scope, ".", function () {
            $("#event-detail-more-button").first().dropdown('toggle');
        });

    };

    return directive;
});

app.directive("keyTable", function () {

    directive = {
        restrict: "A"
    };

    directive.scope = {
        rows: "=keyTableRows",
        activeRowIndex: "=keyTableActiveRowIndex"
    };

    directive.controller = function ($scope, Keyboard, Util, $element) {

        console.log("keyTable");

        keyTableScope = $scope;

        $scope.$element = $element;
        $scope.Keyboard = Keyboard;
        $scope.activeRowIndex = 0;

        var scrollToView = function () {

            var rowIndexClass = "row-index-" + $scope.activeRowIndex;
            var row = angular.element($element).find("." + rowIndexClass);
            if (row.hasClass(rowIndexClass)) {
                Util.scrollElementIntoView(row);
            }
            else {
                Util.scrollElementIntoView(
                    angular.element(
                        $element).find("tr").eq($scope.activeRowIndex));
            }
        };

        Keyboard.scopeBind($scope, "j", function () {
            $scope.$apply(function () {
                if ($scope.activeRowIndex < $scope.rows.length - 1) {
                    $scope.activeRowIndex++;
                }
                scrollToView();
            });
        });

        Keyboard.scopeBind($scope, "k", function () {
            $scope.$apply(function () {
                if ($scope.activeRowIndex > 0) {
                    $scope.activeRowIndex--;
                }
                scrollToView();
            });
        });

        Keyboard.scopeBind($scope, "H", function (e) {
            $scope.$apply(function () {
                $(window).scrollTop(0);
                $scope.activeRowIndex = 0;
            });
        });

        Keyboard.scopeBind($scope, "G", function (e) {
            $scope.$apply(function () {
                $(window).scrollTop($(document).height())
                $scope.activeRowIndex = $scope.rows.length - 1;
            });
        });

    };

    return directive;
});