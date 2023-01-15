/*
 * Licensed to Jasig under one or more contributor license
 * agreements. See the NOTICE file distributed with this work
 * for additional information regarding copyright ownership.
 * Jasig licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License.  You may obtain a
 * copy of the License at the following location:
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
function changeLoginVis() {
    var unobj = document.getElementById("username");
    var unvalue = unobj.value;

    var obj = document.createElement('input');
    obj.id = "username";
    obj.name = "username";
    obj.className = 'field form-control';
    obj.setAttribute('tabindex', '1');
    obj.setAttribute('accesskey', 'i');

    if (document.getElementById("fm1").loginVis.checked) {
        obj.setAttribute('type', 'password');
        unobj.parentNode.replaceChild(obj, unobj);
    } else {
        obj.setAttribute('type', 'text');
        unobj.parentNode.replaceChild(obj, unobj);
    }

    obj.value = unvalue;

    return true;
}

$(window).resize(function(){
    $('#body').css({
        position:'absolute',
        left: ($(window).width() - $('#body').outerWidth())/2,
        top: ($(window).height() - $('#body').outerHeight())/2
    });
});


$(document).ready(function(){
    $("#loginVis").click({id: "username"}, changeLoginVis);
    //focus username field
    $("input:visible:enabled:first").focus();
    //flash error box
    $('#msg.errors').animate({ backgroundColor: 'rgb(187,0,0)' }, 30).animate({ backgroundColor: 'rgb(255,238,221)' }, 500);

    //flash success box
    $('#msg.success').animate({ backgroundColor: 'rgb(51,204,0)' }, 30).animate({ backgroundColor: 'rgb(221,255,170)' }, 500);

    if ( !window.console ) {
        window.console = new Array();
    }
    if (!window.console || window.console == {}) {
        window.console.log = function() {};
    }
    // center page
    $(window).resize();
});
