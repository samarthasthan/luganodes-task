"use client";

import { useState, useEffect } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

// Function to convert UNIX timestamp to human-readable time
const formatTimestamp = (timestamp) => {
  const date = new Date(timestamp * 1000); // Assuming timestamp is in seconds
  return date.toLocaleString();
};

// Function to convert Wei to Ether
const formatFee = (wei) => {
  const etherValue = parseFloat(wei) / 10 ** 18; // Convert Wei to Ether
  return etherValue.toFixed(6); // Display up to 6 decimal places
};

export default function Home() {
  const [deposit, setDeposit] = useState([]);
  const [loading, setLoading] = useState(true); // Loader state

  // Function to fetch data
  const fetchDeposits = () => {
    setLoading(true);
    fetch(`/api/deposits`)
      .then((response) => response.json())
      .then((data) => {
        setDeposit(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Error fetching data:", error);
        setLoading(false);
      });
  };

  // Fetch data on component mount
  useEffect(() => {
    fetchDeposits();
  }, []);

  return (
    <div className="p-8">
      <div className="py-5">
        <h1 className="font-bold text-4xl">Welcome to Luganodes Task</h1>
        <h2 className="font-bold text-2xl">Made by Samarth Asthan 21BRS1248</h2>
      </div>

      {/* Loader Spinner */}
      {loading ? (
        <div className="flex justify-center items-center">
          <div className="loader" />
          <p>Loading...</p>
        </div>
      ) : (
        <>
          {/* Displaying Cards */}
          {deposit.length > 0 ? (
            deposit.map((item, index) => (
              <Card className="my-4" key={index}>
                <CardHeader>
                  <CardTitle>ID: {item.ID}</CardTitle>
                  <CardDescription>
                    Blocknumber: {item.Blocknumber}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <p>
                    <p className="font-bold">Blocktimestamp: </p>
                    {formatTimestamp(item.Blocktimestamp)}
                  </p>
                  <p>
                    <p className="font-bold">Fee:</p> {formatFee(item.Fee)} ETH
                  </p>
                  <p style={{ wordWrap: "break-word" }}>
                    <p className="font-bold">Hash:</p> {item.Hash}
                  </p>
                  <p style={{ wordWrap: "break-word" }}>
                    <p className="font-bold">Pubkey:</p> {item.Pubkey}
                  </p>
                </CardContent>
              </Card>
            ))
          ) : (
            <p>No data available</p>
          )}
        </>
      )}

      {/* Styling for the loader */}
      <style jsx>{`
        .loader {
          border: 4px solid rgba(0, 0, 0, 0.1);
          border-left-color: #09f;
          border-radius: 50%;
          width: 36px;
          height: 36px;
          animation: spin 1s linear infinite;
        }

        @keyframes spin {
          to {
            transform: rotate(360deg);
          }
        }

        /* Ensure hash and public key wrap on small screens */
        @media (max-width: 600px) {
          p {
            word-wrap: break-word;
          }
        }
      `}</style>
    </div>
  );
}
