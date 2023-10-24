/*
Stix2viz and d3 are packaged in a way that makes them work as Jupyter
notebook extensions.  Part of the extension installation process involves
copying them to a different location, where they're available via a special
"nbextensions" path.  This path is hard-coded into their "require" module
IDs.  Perhaps it's better to use abstract names, and add special config
in all cases to map the IDs to real paths, thus keeping the modules free
of usage-specific hard-codings.  But packaging in a way I know works in
Jupyter (an already complicated environment), and having only this config
here, seemed simpler.  At least, for now.  Maybe later someone can structure
these modules and apps in a better way.
*/
require.config({
  paths: {
    "nbextensions/stix2viz/d3": "/assets/js/d3"
  }
});

define(["domReady!", "stix2viz"], function (document, stix2viz) {


  // Init some stuff
  // For optimization purposes, look into moving these to local variables

  var visualizer;
  selectedContainer = document.getElementById('selection');
  // uploader = document.getElementById('uploader');
  canvasContainer = document.getElementById('canvas-container');
  canvas = document.getElementById('canvas');
  //styles = window.getComputedStyle(uploader);

  /* ******************************************************
   * Will be called right before the graph is built.
   * ******************************************************/
  function vizCallback() {
    hideMessages();
    //resizeCanvas();
  }

  /* ******************************************************
   * Will be called if there's a problem parsing input.
   * ******************************************************/
  function errorCallback() {
    document.getElementById('chosen-files').innerText = "";
    document.getElementById("files").value = "";
  }

  /* ******************************************************
   * Initializes the graph, then renders it.
   * ******************************************************/
  function vizStixWrapper(content) {
    cfg = {
      iconDir: "/assets/img/icons"
    }
    visualizer = new stix2viz.Viz(canvas, cfg, populateLegend, populateSelected);
    visualizer.vizStix(content, '', vizCallback, errorCallback);
  }

  /* ******************************************************
   * Handles content pasted to the text area.
   * ******************************************************/
  function handleTextarea() {
    customConfig = '';
    //customConfig = document.getElementById('paste-area-custom-config').value;
    content = document.getElementById('paste-area-stix-json').value;
    vizStixWrapper(content, customConfig);
  }

  /* ******************************************************
   * Adds icons and information to the legend.
   *
   * Takes an array of type names as input
   * ******************************************************/
  function populateLegend(typeGroups) {
    var ul = document.getElementById('legend-content');
    var color = d3.scale.category20();
    typeGroups.forEach(function (typeName, index) {
      var li = document.createElement('li');
      var val = document.createElement('p');
      var key = document.createElement('div');
      var keyImg = document.createElement('img');
      keyImg.onerror = function () {
        // set the node's icon to the default if this image could not load
        this.src = visualizer.d3Config.iconDir + "/stix2_custom_object_icon_tiny_round_v1.svg";
      }
      keyImg.src = visualizer.iconFor(typeName);
      keyImg.width = "37";
      keyImg.height = "37";
      keyImg.style.background = "radial-gradient(" + color(index) + " 16px,transparent 16px)";
      key.appendChild(keyImg);
      val.innerText = typeName.charAt(0).toUpperCase() + typeName.substr(1).toLowerCase(); // Capitalize it
      li.appendChild(key);
      li.appendChild(val);
      ul.appendChild(li);
    });
  }

  /* ******************************************************
   * Adds information to the selected node table.
   *
   * Takes datum as input
   * ******************************************************/
  function populateSelected(d) {
    // Remove old values from HTML
    selectedContainer.innerHTML = "";

    var counter = 0;

    Object.keys(d).forEach(function (key) { // Make new HTML elements and display them
      // Create new, empty HTML elements to be filled and injected
      var div = document.createElement('div');
      var type = document.createElement('div');
      var val = document.createElement('div');

      // Assign classes for proper styling
      if ((counter % 2) != 0) {
        div.classList.add("odd"); // every other row will have a grey background
      }
      type.classList.add("type");
      val.classList.add("value");

      // Add the text to the new inner html elements
      var value = d[key];
      type.innerText = key;
      val.innerText = value;

      // Add new divs to "Selected Node"
      div.appendChild(type);
      div.appendChild(val);
      selectedContainer.appendChild(div);

      // increment the class counter
      counter += 1;
    });
  }

  /* ******************************************************
   * Hides the data entry container and displays the graph
   * container
   * ******************************************************/
  function hideMessages() {
    //console.log($( "#canvas-container" ).hasClass("hidden"))
    if ($( "#canvas-container" ).hasClass("hidden")) {
      canvasContainer.classList.toggle("hidden");
    }
  }

  /* ******************************************************
   * Resets the graph so new data can be added
   * ******************************************************/
  function reset() {
    $('#legend-content').empty();
    $('#selection').empty();
    visualizer.vizReset();
  }

    /* ******************************************************
   * Checks if visualizer has been created 
   * ******************************************************/
  function getVisualizer(){
    return visualizer != null;
  }

  //Exports functions so they can be called in another js file
  return {
    vizStixWrapper: vizStixWrapper,
    reset: reset,
    getVisualizer: getVisualizer
  }
});
