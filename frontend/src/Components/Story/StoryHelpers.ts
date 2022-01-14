import blob1 from "../../assets/blobs/blob1.png";
import blob2 from "../../assets/blobs/blob2.png";
import blob3 from "../../assets/blobs/blob3.png";
import blob4 from "../../assets/blobs/blob4.png";
import blob5 from "../../assets/blobs/blob5.png";
import blob6 from "../../assets/blobs/blob6.png";
import blob7 from "../../assets/blobs/blob7.png";
import blob8 from "../../assets/blobs/blob8.png";

export default function GetBlob(agentType: string) {
  switch (agentType) {
    case "Team1Agent1":
      return blob1;
    case "Team2":
      return blob2;
    case "Team3":
      return blob3;
    case "Team4":
      return blob4;
    case "Team5":
      return blob5;
    case "Team6":
      return blob6;
    case "Team7":
      return blob7;
    default:
      return blob8;
  }
}
