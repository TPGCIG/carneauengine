export default function SuccessPage() {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen">
        <h1 className="text-4xl font-bold text-green-600 mb-4">Payment Successful!</h1>
        <p className="text-lg text-gray-700">Thank you for your purchase. A confirmation email has been sent.</p>
        <a href="/" className="mt-8 text-blue-500 hover:underline">Return Home</a>
      </div>
    );
  }
