export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-center font-mono text-sm">
        <h1 className="text-4xl font-bold text-center mb-8">
          Language Learning Platform
        </h1>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
          <div className="p-6 border rounded-lg hover:shadow-lg transition">
            <h2 className="text-2xl font-semibold mb-2">For Students</h2>
            <p className="text-gray-600">
              Browse courses, enroll, and learn from expert teachers worldwide
            </p>
          </div>

          <div className="p-6 border rounded-lg hover:shadow-lg transition">
            <h2 className="text-2xl font-semibold mb-2">For Teachers</h2>
            <p className="text-gray-600">
              Create courses, manage students, and share your knowledge
            </p>
          </div>

          <div className="p-6 border rounded-lg hover:shadow-lg transition">
            <h2 className="text-2xl font-semibold mb-2">Live Classes</h2>
            <p className="text-gray-600">
              Join live video sessions and interact with teachers in real-time
            </p>
          </div>
        </div>

        <div className="flex gap-4 justify-center mt-12">
          <button className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition">
            Get Started
          </button>
          <button className="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50 transition">
            Browse Courses
          </button>
        </div>
      </div>
    </main>
  )
}
