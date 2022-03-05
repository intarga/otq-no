// This script contains style corrections that could not be performed with CSS only

window.onload = () => {
  fixRainbows();
};

window.addEventListener("resize", () => {
  fixRainbows();
});

// On certain screen sizes/zoom settings, extremity-rainbow shows above/below the main rainbow
// This ensures that both extremities and main rainbow always have the same height
const fixRainbows = () => {
  const sampleListItem = document.getElementsByClassName("top-nav-listitem")[0];
  sampleListItem != null && sampleListItem != undefined
    ? (() => {
        const wrapper = document.getElementById("all-encompasser");
        const extremitiesH = window
          .getComputedStyle(wrapper)
          .getPropertyValue("background-size");
        const rainbowUnitH = window
          .getComputedStyle(sampleListItem)
          .getPropertyValue("border-top-width");

        // Replace extremity height if extremity has defined value, and it is different from the main rainbow's border-top width
        if (extremitiesH !== "auto") {
          const extremitiesVals = extremitiesH.split(" ");
          const extremitiesLastVal =
            extremitiesVals[extremitiesVals.length - 1];
          if (extremitiesLastVal !== rainbowUnitH) {
            wrapper.style.backgroundSize = `${
              extremitiesVals[0]
            } ${rainbowUnitH}, ${
              extremitiesVals[2] ? extremitiesVals[2] : extremitiesVals[0]
            } ${rainbowUnitH}`;
          }
        }
      })()
    : "";
};
