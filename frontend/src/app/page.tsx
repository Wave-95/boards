import Footer from '@/components/footer';
import WidthContainer from '@/components/widthContainer';
import { COOKIE_NAME_JWT_TOKEN, NAVBAR_HEIGHT_PX } from '@/constants';
import { cookies } from 'next/headers';
import Link from 'next/link';
import { FaPaperPlane, FaStickyNote, FaRegFileCode } from 'react-icons/fa';

async function getToken() {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  return jwtToken;
}

export default async function Page() {
  const token = await getToken();
  return (
    <div>
      <div className="hero bg-base-200" style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}>
        <div className="hero-content flex-col lg:flex-row lg:space-x-6">
          <div>
            <h1 className="text-5xl font-bold">Boards</h1>
            <p className="py-6 max-w-sm text-xl text-gray-700">
              Collaborate with others and turn your ideas into actions
            </p>
            <div className="flex space-x-4">
              <Link href="/auth/signup" className="btn btn-primary">
                Sign Up
              </Link>
              <Link href={token ? '/dashboard' : '/auth/signin'} className="btn btn-secondary btn-outline">
                Visit app
              </Link>
            </div>
          </div>
          <img src="/Hero.png" style={{ width: '700px' }} className="rounded-lg shadow-2xl" alt="Boards app" />
        </div>
      </div>
      <div className="bg-white">
        <WidthContainer className="h-full flex flex-col justify-evenly items-center my-8">
          <div className="flex flex-col lg:flex-row items-center justify-between w-full my-8">
            <div className="card w-80 border !rounded-none">
              <div className="card-body">
                <FaPaperPlane size={30} />
                <h2 className="card-title">Invite</h2>
                <p className="text-gray-700">Invite other collaborators and work together in real-time</p>
              </div>
            </div>
            <div className="card w-80 border !rounded-none">
              <div className="card-body">
                <FaStickyNote size={30} />
                <h2 className="card-title">Organize</h2>
                <p className="text-gray-700">Capture ideas, tasks, and goals on digital sticky notes</p>
              </div>
            </div>
            <div className="card w-80 border !rounded-none">
              <div className="card-body">
                <FaRegFileCode size={30} />
                <h2 className="card-title">Automate</h2>
                <p className="text-gray-700">Automatically group similar posts together and summarize content</p>
              </div>
            </div>
          </div>
          <div className="w-full flex flex-col lg:flex-row my-8">
            <div className="flex-1 my-4">
              <div>
                <h2 className="font-bold text-xl pb-4">Documentation</h2>
                <p className="max-w-lg text-gray-700">
                  Interested in the backend components? Take a look at the documentation for the REST API or WebSocket
                  events.
                </p>
                <div className="flex space-x-6 pt-6">
                  <a className="btn btn-primary" href={process.env.NEXT_PUBLIC_API_DOCS_URL} target="_blank">
                    API Docs
                  </a>
                  <a className="btn btn-outline" href={process.env.NEXT_PUBLIC_API_DOCS_URL} target="_blank">
                    WebSocket Docs
                  </a>
                </div>
              </div>
            </div>
            <div className="flex-1 my-4">
              <div className="flex flex-col justify-between">
                <h2 className="font-bold text-xl pb-4 pt-8 lg:pt-0">Boards is built with</h2>
                <div className="flex justify-between">
                  <img src="/golang.png" width={50} />
                  <img src="/nextjs.png" width={50} />
                  <img src="/postgresql.png" width={50} />
                  <img src="/redis.png" width={50} />
                </div>
                <div className="flex justify-between mt-4">
                  <img src="/docker.png" width={50} />
                  <img src="/kubernetes.png" width={50} />
                  <img src="/aws.png" width={50} />
                  <img src="/rabbitmq.png" width={50} />
                </div>
              </div>
            </div>
          </div>
        </WidthContainer>
      </div>
      <Footer />
    </div>
  );
}
