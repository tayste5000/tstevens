window.addEventListener('load', setup);
var viewer;
var presentation;
var dialogue;

function setup() {

var options = {
  width: "auto",
  height: 500,
  background: "#fdfdfd",
  antialias: true,
  quality: 'medium',
  outlineWidth: 1.0
};

viewer = pv.Viewer(document.getElementById('viewer'), options);
dialogue = document.querySelector("#dialogue");

viewer.on("viewerReady", function(){

  dialogue.innerHTML += "<button type=\"button\" id=\"start\">Start</button>";

  var startBtn = document.querySelector("#start");

  presentation = makePresentation(dialogue);

  startBtn.addEventListener("click", function(){
    presentation.start();
  });
})

}

function makePresentation(dialogue){

var count;
var reverse = false;
var amino_acids = ["A", "C", "D", "E", "F", "G", "H", "I", "K", "L", "M", "N", "P", "Q", "R", "S", "T", "V", "W", "Y"];
var structures = {};

/*init*/
loadStructure = function loadStructure(filename,structureName){
  pv.io.fetchPdb("/public/pdbs/" + filename, function(output){structures[structureName] = output;});
}

amino_acids.forEach(function(value,index){
  loadStructure("amino_acids/" + value + ".pdb", value);
});

loadStructure("ss/polypeptide.pdb", "polypeptide");
loadStructure("ss/helix.pdb", "helix");
loadStructure("ss/lysozyme-naked.pdb", "lysozyme");
var lysozymeCount;
loadStructure("ss/estrogen.pdb", "estrogen");
loadStructure("ss/estrogen-4HT.pdb", "tamoxifen");

/*return presentation object*/
return {
  "start": function start(){

    dialogue.innerHTML = "<p>Proteins perform many functions in our cells. " +
      "Proteins are made up of small molecules called amino acids, that contain carbon, <span style=\"color:blue\">nitrogen</span>, and <span style=\"color:red\">oxygen</span>." +
    "</p>" +
    "<button type=\"button\" id=\"next\">Next</button>";

    count = 0;
    viewer.setCenter([0, 0, 0]);
    viewer.setZoom(25);
    viewer.hide("*")
    viewer.ballsAndSticks(amino_acids[count],structures[amino_acids[count]]);
    count++;

    next = document.querySelector("#next");

    next.addEventListener("click", function(){
      presentation.aminoAcids.start();
    });
  },
  "aminoAcids": {
    "stop": false,
    "start": function aminoAcidsStart(){
      dialogue.innerHTML = "<p>There are 20 different amino acids, and each one has a different side chain, giving it unique chemical properties.</p>" +
      "<button id=\"prev\">Previous</button>" +
      "<button id=\"next\">Next</button>";

      viewer.setRotation(
        [-0.8845363259315491, 0.44514814019203186, 0.13903594017028809, 0,
        -0.46628516912460327, -0.8389679789543152, -0.2803587019443512, 0,
        -0.008154172450304031, -0.31285709142684937, 0.9497630000114441, 0,
        0, 0, 0, 1],
      300);
      viewer.setCenter([0, 0, 0]);

      var next = document.querySelector("#next");
      var next = document.querySelector("#next");

      next.addEventListener("click", function(){
        presentation.aminoAcids.stop = true;
        window.setTimeout(function(){presentation.polypeptide.start();},300);
      });

      prev.addEventListener("click", function(){
        presentation.aminoAcids.stop = true;
        window.setTimeout(function(){presentation.start();},300);
      });

      presentation.aminoAcids.stop = false;
      count = 1;
      presentation.aminoAcids.loop();
    },
    "loop": function aminoAcidsLoop(){
      if (presentation.aminoAcids.stop == true) return;

      window.setTimeout(function(){ 

        viewer.hide("*");
        viewer.ballsAndSticks(amino_acids[count],structures[amino_acids[count]],{radius: 0.3});

        if (count == amino_acids.length - 1) count = 0;

        else count++;

        aminoAcidsLoop();

      },250);
    }
  },
  "polypeptide": {
    "start": function polypeptideStart(){
      dialogue.innerHTML = "<p>Our cells use amino acids to build polypeptide chains, better known as proteins.</p>" +
      "<button type=\"button\" id=\"prev\">Previous</button>" +
      "<button type=\"button\" id=\"next\">Next</button>";

      var next = document.querySelector("#next");
      var prev = document.querySelector("#prev");

      next.addEventListener("click", function(){
        presentation.polypeptide.stop = true;
        window.setTimeout(function(){presentation.helix.start();},300);
      });

      prev.addEventListener("click", function(){
        presentation.polypeptide.stop = true;
        window.setTimeout(function(){presentation.aminoAcids.start();},300);
      });

      viewer.hide("*");
      viewer.centerOn(structures["polypeptide"], 300);
      viewer.setZoom(30, 300);
      viewer.setRotation([0.884534478187561, 0.4451482594013214, -0.13903702795505524, 0,
        0.4662850499153137, -0.8389653563499451, 0.28037029504776, 0,
        0.008155263029038906, -0.31284546852111816, -0.9497674703598022, 0,
        0, 0, 0, 1], 300);

      presentation.polypeptide.stop = false;
      viewer.ballsAndSticks("polypeptide", structures["polypeptide"].select({"rnumRange":[2,2]}));
      count = 1;

      presentation.polypeptide.loop();

    },
    "stop": false,
    "loop": function polypeptideLoop(){
      if (presentation.polypeptide.stop) return;

      window.setTimeout(function(){
        viewer.hide("*");
        viewer.ballsAndSticks("polypeptide", structures["polypeptide"].select({"rnumRange":[2,count+2]}));

        if (count==3){
          count=0;
          window.setTimeout(function(){polypeptideLoop();},1000);
        }
        else {
          count++;
          polypeptideLoop();
        }
      },100)
    }
  },
  "helix": {
    "start": function helixStart(){
      dialogue.innerHTML = "<p>Proteins can form a variety of shapes, such as this alpha helix.</p>" +
      "<button type=\"button\" id=\"prev\">Previous</button>" +
      "<button type=\"button\" id=\"next\">Next</button>";

      var next = document.querySelector("#next");
      var prev = document.querySelector("#prev");

      next.addEventListener("click", function(){
        presentation.helix.stop = true;
        window.setTimeout(function(){presentation.lysozyme.start();},200);
      });

      prev.addEventListener("click", function(){
        presentation.helix.stop = true;
        window.setTimeout(function(){presentation.polypeptide.start();},200);
      });

      viewer.hide("*");
      viewer.centerOn(structures["helix"], 300);
      viewer.setZoom(40, 300);
      viewer.setRotation(
        [-0.06436467915773392, 0.7501113414764404, 0.6580818295478821, 0,
        -0.35576286911964417, -0.6333650350570679, 0.6871384382247925, 0,
        0.9323446154594421, -0.18990138173103333, 0.30765125155448914, 0,
        0, 0, 0, 1], 300);

      presentation.helix.stop = false;
      viewer.ballsAndSticks("helix", structures["helix"].select({"rnumRange":[22,22]}));
      count = 1;

      presentation.helix.loop();

    },
    "stop": false,
    "loop": function helixLoop(){
      if (presentation.helix.stop) return;

      window.setTimeout(function(){

        if (count == 13) {
          count = 0;
          viewer.hide("*")
          viewer.cartoon("helix", structures["helix"]);
          window.setTimeout(function(){helixLoop();},1000);
        }

        else {
          viewer.hide("*")
          viewer.ballsAndSticks("helix", structures["helix"].select({"rnumRange":[22,count+22]}),{radius: 0.2});
          count++;
          helixLoop();
        }

      },100);
    }
  },
  "lysozyme": {
    "start": function lysozymeStart(){
      dialogue.innerHTML = "<p>Proteins can contain over 1000 amino acids, leading to a nearly infinite number of shapes they can fold into.</p>" +
      "<button type=\"button\" id=\"prev\">Previous</button>" +
      "<button type=\"button\" id=\"next\">Next</button>";

      var next = document.querySelector("#next");
      var prev = document.querySelector("#prev");

      next.addEventListener("click", function(){
        presentation.lysozyme.stop = true;
        window.setTimeout(function(){presentation.estrogen.start();},300);
      });

      prev.addEventListener("click", function(){
        presentation.lysozyme.stop = true;
        window.setTimeout(function(){presentation.helix.start();},300);
      });

      viewer.hide("*");
      viewer.centerOn(structures["lysozyme"], 300);
      viewer.setZoom(100, 300);
      viewer.setRotation([-0.16544285416603088, 0.04722153767943382, 0.9850873351097107, 0,
        -0.04623271897435188, 0.9973832964897156, -0.05557561293244362, 0,
        -0.985134482383728, -0.054737888276576996, -0.16282692551612854, 0,
        0, 0, 0, 1], 300);

      lysozymeCount = structures["lysozyme"].residueCount();

      presentation.lysozyme.stop = false;
      viewer.ballsAndSticks("lysozyme", structures["lysozyme"].select({"rnumRange":[1,1]}));
      count = 1;

      presentation.lysozyme.loop();
    },
    "stop": false,
    "loop": function lysozymeLoop(){
      if (presentation.lysozyme.stop) return;

      window.setTimeout(function(){

        if (count >= lysozymeCount + 1) {
          viewer.hide("*")
          viewer.cartoon("lysozyme", structures["lysozyme"]);
        }

        else {
          viewer.hide("*")
          viewer.ballsAndSticks("lysozyme", structures["lysozyme"].select({"rnumRange":[1,count+1]}),{radius: 0.2});
          count += 10;
          lysozymeLoop();
        }

      },250)
    }
  },
  "estrogen": {
    "start": function estrogenStart(){
      dialogue.innerHTML = "<p>This is the estrogen receptor. Overactivation of the estrogen receptor can lead to breast cancer.</p>" +
      "<button type=\"button\" id=\"prev\">Previous</button>" +
      "<button type=\"button\" id=\"next\">Next</button>";

      var next = document.querySelector("#next");
      var prev = document.querySelector("#prev");

      next.addEventListener("click", function(){
        presentation.estrogen.pocket();
      });
      prev.addEventListener("click", function(){
        presentation.lysozyme.start();
      });

      viewer.hide("*");
      viewer.centerOn(structures["estrogen"], 300);
      viewer.setZoom(100, 300);
      viewer.setRotation([-0.6491602063179016, -0.692853569984436, 0.31387078762054443, 0,
        -0.7605252265930176, 0.5982341170310974, -0.25238102674484253, 0,
        -0.012906712479889393, -0.40254804491996765, -0.9152800440788269, 0,
        0, 0, 0, 1], 300);

      viewer.cartoon("estrogen", structures["estrogen"]);
    },
    "pocket": function estrogenPocket(){
      dialogue.innerHTML = "<p>The estrogen receptor structure contains a binding pocket in the shape of estrogen that is lined with chemical groups attracted to estrogen.</p>" +
      "<button type=\"button\" id=\"prev\">Previous</button>" +
      "<button type=\"button\" id=\"next\">Next</button>";

      var next = document.querySelector("#next");
      var prev = document.querySelector("#prev");

      next.addEventListener("click", function(){
        presentation.estrogen.tamoxifen();
      });

      prev.addEventListener("click", function(){
        presentation.estrogen.start();
      });
      
      viewer.hide("*")
      viewer.cartoon("estrogen", structures["estrogen"])
      window.setTimeout(function(){viewer.ballsAndSticks("estrogen", structures["estrogen"].select({"rname":"EST"}), {"color":pv.color.uniform("yellow")});},500);

      viewer.setRotation([-0.4869979918003082, -0.7540156841278076, 0.4407404363155365, 0,
        -0.8724094033241272, 0.44379809498786926, -0.20472747087478638, 0,
        -0.04123280569911003, -0.4842182993888855, -0.8739356994628906, 0,
        0, 0, 0, 1], 300);
      viewer.setZoom(50,300);

    },
    "tamoxifen": function estrogenStart(){
      dialogue.innerHTML = "<p>Tamoxifin is a common treatment for hormone related breast cancer. It works by binding to the estrogen binding pocket and forcing the receptor into an inactive shape.</p>" +
      "<button type=\"button\" id=\"prev\">Previous</button>";

      var prev = document.querySelector("#prev");

      prev.addEventListener("click", function(){
        presentation.estrogen.pocket();
      });

      viewer.setRotation([-0.21546663343906403, -0.8086357712745667, 0.5473734736442566, 0,
        -0.8190129399299622, 0.45492640137672424, 0.3496115505695343, 0,
        -0.531719446182251, -0.3729931712150574, -0.7603074312210083, 0,
        0, 0, 0, 1], 300);

      viewer.setZoom(100,300);
      window.setTimeout(function(){
        viewer.hide("*");
        viewer.cartoon("tamoxifen", structures["tamoxifen"]);
      },500);

      window.setTimeout(function(){
        viewer.ballsAndSticks("tamoxifin", structures["tamoxifen"].select({"rname":"OHT"}),{"color":pv.color.uniform("yellow")});
      },500)
    }
  }
};
}