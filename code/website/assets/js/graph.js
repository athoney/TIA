require(["application2", "piechart", "barchart"], function (app2, piechart, barchart) {
    /* Graphs under visualization tab*/
    var pieChart2;
    var pieChart3;
    var barChart;
    /*
    * Runs when a new Vendor is selected
    * Sends a post request with the vendor
    * Populates the product dropdown accordingly
    */
    $('#vendorSelect').change(function () {
        //Destroy graphs that already exist
        destroy();
        //stores value of tag with id = 'vendorSelect'
        var vendor = $('#vendorSelect').val();
        //Sets <span id="vendorSpan"> to vendor
        $('#vendorSpan').text(vendor);
        $('#productSpan').text("");
        //Empties option tags in <select id='productSelect'> tag
        $('#productSelect').find('option').remove().end();
        //Removes disabled attribute
        $('#productSelect').removeAttr('disabled');
        /*
        * Sends post request to the server at graph/(vendor name)
        * URL matches POST request in terminal logging
        */
        var posting = $.get("/graph/" + vendor);
        //When the post is finished:
        posting.done(function (data) {
            //Shove returned data (which is returned as html) in <select id='productSelect'> tag
            $('#productSelect').html(data);
        });
        if (!$("#canvas-container").hasClass("hidden")) {
            $("#canvas-container").addClass("hidden");
        }
    });

    /*
    * Runs when a new Product is selected
    * Sends a post request with the vendor and product
    * Populates the vulnerabilites table accordingly
    */
    $('#productSelect').change(function () {
        //Destroy graphs that already exist
        destroy();
        //stores value of tag with id = "productSelect"
        var product = $('#productSelect').val();
        //Sets <span id="productSpan"> to product
        $('#productSpan').text(product);
        //stores value of tag with id = 'vendorSelect'
        var vendor = $('#vendorSelect').val();
        /*
        * Sends post request to the server at graph/(vendor name)/(product name)
        * URL matches POST request in terminal logging
        */
        var posting = $.get("/graph/" + vendor + "/" + product);
        //When the post is finished:
        posting.done(function (data) {
            //Shove returned data (which is returned as html) in <tbody id='vulnTable'> tag
            $('#vulnTable').html(data);
            //Triggers vulnTable on change event
            $('#vulnTable').trigger("change");
        });
    });

    /*
    * Runs when a vendor and product are selected
    * Sends a post request with the vendor and product
    * Renders STIX graph
    * Calls functions to add charts seen in Visualizations tab
    */
    $('#vulnTable').change(function () {
        vendor = $('#vendorSelect').val();
        product = $('#productSelect').val();
        var posting = $.get("/graph/" + vendor + "/" + product + "/stix");
        posting.done(function (data) {
            //Calls function to create Stix graph, legend, and selected node 
            app2.vizStixWrapper(data.stix);
            //Add href to button and set filename for user to download
            $('#download-btn').attr("href", "assets/resources-tmp/" + data.filename)
            $('#download-btn').attr("download", data.filename)
            //Functions to add charts
            addBar();
            addPie();
        })
    });

    /*
    * Runs when a new tab is pressed
    * Rerenders the bar and piecharts if
    * vendor and product are set
    */
    $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
        vendor = $('#vendorSelect').val();
        product = $('#productSelect').val();
        if (!(vendor == null || product == null)) {
            addBar();
            addPie();
        }
    });

    /*
    * Runs when a new cvss tab is pressed
    * Rerenders the piecharts if
    * vendor and product are set
    */
    $('button[data-toggle="tab"]').on('shown.bs.tab', function (e) {
        vendor = $('#vendorSelect').val();
        product = $('#productSelect').val();
        if (!(vendor == null || product == null)) {
            addPie();
        }
    });

    /*
    * Resets any drawn graphs
    */
    function destroy() {
        if (pieChart2 != null) {
            pieChart2.resetPie();
            pieChart3.resetPie();
            barChart.resetBar();
        }
        if (app2.getVisualizer()) {
            app2.reset();
        }
    }

    /*
    * Runs once stix data is returned
    * Sends a post request with the vendor and product
    * Renders bar chart
    */
    function addBar() {
        //clear graph if it already exists
        if (barChart != null) {
            barChart.resetBar();
        }
        vendor = $('#vendorSelect').val();
        product = $('#productSelect').val();

        //GET request to /graph/(vendor)/(product)/cwes
        var posting = $.get("/graph/" + vendor + "/" + product + "/cwes");
        cweIds = [];
        cwenames = [];
        counts = [];
        posting.done(function (data) {
            data.cwes.forEach(weakness => {
                if (weakness.CWEID != "NVD-CWE-noinfo") {
                    cweIds.push(weakness.CWEID);
                    cwenames.push(weakness.CWEID + ":\n" + weakness.Name);
                    counts.push(weakness.Count);
                }
            });
            //If cweIds is not empty, create barChart with returned data
            if (cweIds.length != 0) {
                barChart = new barchart.Bar(cweIds, counts, "barChart", cwenames);
            }
        });
    }

    /*
   * Runs once stix data is returned
   * Sends a post request with the vendor and product
   * Renders pie charts
   */
    function addPie() {
        //clear graphs if they already exist
        if (pieChart2 != null) {
            pieChart2.resetPie();
            pieChart3.resetPie();
        }
        vendor = $('#vendorSelect').val();
        product = $('#productSelect').val();

        //GET request to /graph/(vendor)/(product)/scores
        var posting = $.get("/graph/" + vendor + "/" + product + "/scores");
        posting.done(function (data) {
            //Create piecharts with returned data
            pieChart2 = new piechart.Pie(data.cvss2, "pieChart2");
            pieChart3 = new piechart.Pie(data.cvss3, "pieChart3");
        });

    }

    // A funny
    $('button#pressMe').click(function () {
        alert("You want to hear a JavaScript joke?\n\nI'll callback later.");
    });

    $("[data-toggle=tooltip").tooltip();
});
