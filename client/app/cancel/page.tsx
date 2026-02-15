export default function CancelPage() {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen">
        <h1 className="text-4xl font-bold text-red-600 mb-4">Payment Cancelled</h1>
        <p className="text-lg text-gray-700">You have not been charged.</p>
        <a href="/" className="mt-8 text-blue-500 hover:underline">Return Home</a>
      </div>
    );
  }
