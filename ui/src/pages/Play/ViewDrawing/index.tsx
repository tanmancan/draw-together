import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

function ViewDrawing() {
  const { data } = useParams();
  const [imgData, setImgData] = useState("");

  useEffect(() => {
    if (data) {
      const dParsed = atob(data);
      setImgData(dParsed);
    }
  }, [data]);
  return (
    <div>
      <img width={512} height={512} src={imgData} />
    </div>
  );
}

export default ViewDrawing;
