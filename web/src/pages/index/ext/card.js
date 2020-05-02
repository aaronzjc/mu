import MText from "../components/cards/MText";
import MRichText from "../components/cards/MRichText";

const CardMap = {
    0: MText.name,
    1: MRichText.name,
};

const Cards = {
    [MText.name]: MText,
    [MRichText.name]: MRichText
};

export {CardMap, Cards}