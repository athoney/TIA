// Circle following cursor effect
var cursor = document.getElementById("cursor");
document.addEventListener("mousemove", function(e) {
  cursor.style.left = (e.pageX - 25) + "px";
  cursor.style.top = (e.pageY - 25) + "px";
});
