import * as React from "react";
import Helmet from "react-helmet";
import { withRouter, Link } from "react-router-dom";
import { graphql, compose, withApollo } from "react-apollo";

import { listSupportBundles } from "../../queries/TroubleshootQueries";
// import { archiveSupportBundle } from "../../mutations/SupportBundleMutations";

import Loader from "../shared/Loader";
import SupportBundleRow from "./SupportBundleRow";
import "../../scss/components/troubleshoot/SupportBundleList.scss";

class SupportBundleList extends React.Component {

  render() {
    const { watch } = this.props;
    const { loading, error, listSupportBundles } = this.props.listSupportBundles;

    if (error) {
      <p>{error.message}</p>
    }

    let bundlesNode;
    if (listSupportBundles?.length) {
      bundlesNode = (
        listSupportBundles.map(bundle => (
          <SupportBundleRow
            key={bundle.id}
            bundle={bundle}
            watchSlug={this.props.watch.slug}
          />
        ))
      );
    } else {
      bundlesNode = (
        <div className="flex1 flex-column justifyContent--center alignItems--center">
          <div className="flex-column u-textAlign--center alignItems--center">
            <p className="u-fontSize--largest u-color--tundora u-lineHeight--normal u-fontWeight--bold">You haven't generated any support bundles</p>
            <p className="u-marginTop--normal u-fontSize--normal u-color--dustyGray u-fontWeight--normal">Generating bundles is simple and we'll walk you through it, <span onClick={this.toggleModal} className="u-color--astral u-fontWeight--medium u-textDecoration--underlineOnHover">get started now</span></p>
          </div>
        </div>
      );
    }

    
    return (
      <div className="container u-paddingBottom--30 u-paddingTop--30 flex1 flex">
        <Helmet>
          <title>{`${watch.watchName} Troubleshoot`}</title>
        </Helmet>
        <div className="flex1 flex-column">
          <div className="flex flex1">
            <div className="flex1 flex-column">
              <div className="u-position--relative flex-auto u-paddingBottom--10 flex">
                <div className="flex flex1">
                  <div className="flex1 u-flexTabletReflow">
                    <div className="flex flex1">
                      <div className="flex-auto alignSelf--center">
                        <h2 className="u-fontSize--larger u-fontWeight--bold u-color--tuna flex alignContent--center">Support bundles</h2>
                      </div>
                    </div>
                    <div className="RightNode flex-auto flex-column flex-verticalCenter u-position--relative">
                      <Link to={`${this.props.match.url}/generate`} className="btn secondary">Generate a support bundle</Link>
                    </div>
                  </div>
                </div>
              </div>
              <div className="flex1 flex-column u-overflow--auto">
                {loading ?
                  <div className="flex1 flex-column justifyContent--center alignItems--center">
                    <Loader size="60" color="#44bb66" />
                  </div>
                :
                  bundlesNode
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

export default withRouter(compose(
  withApollo,
  graphql(listSupportBundles, {
    name: "listSupportBundles",
    options: props => {
      return {
        variables: {
          watchSlug: props.watch.slug
        },
        fetchPolicy: "no-cache",
      }
    }
  }),
  // graphql(archiveSupportBundle, {
  //   props: ({ mutate }) => ({
  //     archiveSupportBundle: (id) => mutate({ variables: { id } })
  //   })
  // }),
)(SupportBundleList));
