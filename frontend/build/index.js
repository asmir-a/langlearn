const labels = document.getElementsByTagName("h2");
const buttons = document.getElementsByTagName("button");

const handleOnClick = async () => {
  const fetchRandomNumber = async () => {
    console.log("sending request to the api endpoint");
    const response = await fetch("/api/random-number");
    const jsonData = await response.text();
    console.log(`got ${jsonData} from the endpoint`);
    return parseInt(jsonData);
  }
  const updateLabel = (number) => {
    labels[0].innerHTML = `${number}`
  }

  updateLabel(await fetchRandomNumber())
};

buttons[0].addEventListener("click", handleOnClick);
