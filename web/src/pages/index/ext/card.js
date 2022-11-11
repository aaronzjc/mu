import MText from "../components/cards/MText";
import MRichText from "../components/cards/MRichText";
import MVideo from "../components/cards/MVideo";

const CardMap = {
    0: MText.name,
    1: MRichText.name,
    2: MVideo.name,
};

const Cards = {
    [MText.name]: MText,
    [MRichText.name]: MRichText,
    [MVideo.name]: MVideo,
};

export {CardMap, Cards}